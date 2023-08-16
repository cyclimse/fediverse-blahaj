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

type CrawlResult struct {
	Domain   string
	Peers    []string
	Nodeinfo nodeinfo.Nodeinfo
	Err      CrawlError
}

func newResultFromError(domain string, err CrawlError) CrawlResult {
	return CrawlResult{
		Domain: domain,
		Err:    err,
	}
}

func (c *Crawler) Crawl(ctx context.Context, domain string) CrawlResult {
	slog.InfoCtx(ctx, "crawling", "domain", domain)

	// lookup the domain via DNS
	// avoids retrying on such hosts

	_, err := net.LookupIP(domain)
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			err = errCrawlTimeout
		} else {
			err = errDomainLookupFailed.Wrap(err)
		}
		return newResultFromError(domain, err.(CrawlError))
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
			if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
				return newResultFromError(domain, errCrawlTimeout)
			}
		}
		if !canProceed {
			return newResultFromError(domain, errRobotsTxtDisallowsCrawling)
		} else {
			break
		}
	}
	if err != nil {
		// an error occurred, but we can proceed as if the robots.txt allowed crawling
		slog.ErrorCtx(ctx, "failed to acknowledge robots.txt", "error", err)
	}

	// retry for the rest of the requests
	c.client.RetryMax = 3

	nodeInfo, err := c.getNodeInfo(ctx, url)
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			err = errCrawlTimeout
		}
		return newResultFromError(domain, err.(CrawlError))
	}

	peers, err := c.GetPeers(ctx, url, nodeInfo)
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			err = errCrawlTimeout
		}
		return newResultFromError(domain, err.(CrawlError))
	}

	return CrawlResult{
		Domain:   domain,
		Peers:    peers,
		Nodeinfo: nodeInfo,
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
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return true, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return true, newCrawlInternalError(err)
	}

	group := robots.FindGroup(c.userAgent)

	// TODO: fill in more endpoints
	canProceed :=
		group.Test(url+"/.well-known/nodeinfo") &&
			group.Test(url+"/api/v1/instance/peers")

	return canProceed, nil
}

func (c *Crawler) getNodeInfo(ctx context.Context, url string) (nodeinfo.Nodeinfo, CrawlError) {
	r, err := retryablehttp.NewRequest("GET", url+"/.well-known/nodeinfo", nil)
	if err != nil {
		return nil, newCrawlInternalError(err)
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, errNetworkError.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errNetworkError.Wrap(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	var w nodeinfo.WellKnown
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return nil, errNodeInfoSyntax.Wrap(err)
	}

	if len(w.Links) == 0 {
		return nil, errNodeInfoSyntax.Wrap(fmt.Errorf("no links in nodeinfo"))
	}

	link, nodeInfo, err := nodeinfo.HighestSupported(w)
	if err != nil {
		return nil, errNodeInfoSyntax.Wrap(err)
	}

	r, err = retryablehttp.NewRequest("GET", link, nil)
	if err != nil {
		return nil, newCrawlInternalError(err)
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err = c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, errNetworkError.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errNetworkError.Wrap(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, newCrawlInternalError(err)
	}

	err = json.Unmarshal(b, &nodeInfo)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to decode nodeinfo", "error", err, "body", b)
		return nil, errNodeInfoSyntax.Wrap(err)
	}

	return nodeInfo, nil
}
