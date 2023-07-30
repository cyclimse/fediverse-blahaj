package crawler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"golang.org/x/exp/slog"

	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/temoto/robotstxt"
)

func New() *Crawler {
	return &Crawler{
		client:    retryablehttp.NewClient(),
		userAgent: "fediverse-blahaj/0.0.1",
	}
}

type Crawler struct {
	client    *retryablehttp.Client
	userAgent string
}

func (c *Crawler) Crawl(ctx context.Context, domain string) CrawlResult {
	slog.InfoCtx(ctx, "crawling", "domain", domain)

	// lookup the domain via DNS
	// avoids retrying on such hosts

	_, err := net.LookupIP(domain)
	if err != nil {
		return resultFromError(domain, fmt.Errorf("failed to lookup domain: %w", err))
	}

	// do not retry for the first request
	// for example, if the port is not open, we will get a connection refused error
	// and we don't want to retry that
	c.client.RetryMax = 0

	var url string
	for _, prefix := range []string{"https://", "http://"} {
		url = prefix + domain
		canProceed, err := c.acknowledgeRobotsTxt(ctx, url)
		if err != nil {
			if errors.Is(err, &tls.CertificateVerificationError{}) {
				// will use http instead
				continue
			}
			// if the port is not open, we will get a connection refused error
			slog.ErrorCtx(ctx, "failed to acknowledge robots.txt", "error", err)
		}
		if !canProceed {
			return resultFromError(domain, fmt.Errorf("robots.txt disallows crawling"))
		} else {
			break
		}
	}

	// retry for the rest of the requests
	c.client.RetryMax = 3

	nodeInfo, err := c.getNodeInfo(ctx, url)
	if err != nil {
		return resultFromError(domain, fmt.Errorf("failed to get nodeinfo: %w", err))
	}

	peers, err := c.GetPeers(ctx, url, nodeInfo)
	if err != nil {
		return resultFromError(domain, fmt.Errorf("failed to get peers: %w", err))
	}

	return CrawlResult{
		Domain:   domain,
		Peers:    peers,
		NodeInfo: nodeInfo,
	}
}

// AcknowledgeRobotsTxt checks if the robots.txt allows crawling the given node.
// Assumes true unless the robots.txt explicitly disallows crawling.
func (c *Crawler) acknowledgeRobotsTxt(ctx context.Context, url string) (bool, error) {
	r, err := retryablehttp.NewRequest("GET", url+"/robots.txt", nil)
	if err != nil {
		return true, err
	}

	r.Header.Set("User-Agent", c.userAgent)

	c.client.RetryMax = 0
	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return true, fmt.Errorf("failed to get robots.txt: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("failed to get robots.txt, status code: %d", resp.StatusCode)
	}

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return true, fmt.Errorf("failed to parse robots.txt: %w", err)
	}

	group := robots.FindGroup(c.userAgent)

	// TODO: fill in more endpoints
	canProceed :=
		group.Test(url+"/.well-known/nodeinfo") &&
			group.Test(url+"/api/v1/instance/peers")

	return canProceed, nil
}

func (c *Crawler) getNodeInfo(ctx context.Context, url string) (nodeinfo.Nodeinfo, error) {
	r, err := retryablehttp.NewRequest("GET", url+"/.well-known/nodeinfo", nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("getting well-known: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting well-known, status code: %d", resp.StatusCode)
	}

	var w nodeinfo.WellKnown
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("decoding well-known: %w", err)
	}

	if len(w.Links) == 0 {
		return nil, fmt.Errorf("no nodeinfo link found")
	}

	link, nodeInfo, err := nodeinfo.HighestSupported(w)
	if err != nil {
		return nil, fmt.Errorf("finding highest supported nodeinfo version: %w", err)
	}

	r, err = retryablehttp.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err = c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("getting nodeinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting nodeinfo, status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading nodeinfo: %w", err)
	}

	err = json.Unmarshal(b, &nodeInfo)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to decode nodeinfo", "error", err, "body", b)
		return nil, fmt.Errorf("decoding nodeinfo: %w", err)
	}

	return nodeInfo, nil
}
