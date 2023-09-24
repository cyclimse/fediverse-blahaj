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
	"net/netip"
	"time"

	"log/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/temoto/robotstxt"
)

func New(userAgent string) *Crawler {
	return &Crawler{
		client:    retryablehttp.NewClient(),
		userAgent: userAgent,
	}
}

type Crawler struct {
	client    *retryablehttp.Client
	userAgent string
}

type CrawlResult struct {
	slog.LogValuer

	Domain  string
	Err     error
	ErrCode models.CrawlErrCode

	Start       time.Time
	End         time.Time
	ResolvedIPs []net.IP
	RawNodeinfo json.RawMessage
	Nodeinfo    nodeinfo.Nodeinfo
	Peers       []string
}

func CrawlFromResult(r CrawlResult) models.Crawl {
	addresses := make([]netip.Addr, 0, len(r.ResolvedIPs))
	for _, ip := range r.ResolvedIPs {
		addr, ok := netip.AddrFromSlice([]byte(ip))
		if !ok {
			slog.Error("failed to convert ip to netip.Addr", "ip", ip)
			continue
		}
		addresses = append(addresses, addr)
	}

	c := models.Crawl{
		Domain:    r.Domain,
		Addresses: addresses,

		StartedAt:  r.Start,
		FinishedAt: r.End,
		Status:     models.CrawlStatusCompleted,

		Peers:         r.Peers,
		NumberOfPeers: new(int32),

		SoftwareName:    new(string),
		SoftwareVersion: new(string),

		OpenRegistrations: new(bool),
	}

	if r.Err != nil {
		c.Status = models.CrawlStatusFailed
		c.Err = &models.CrawlError{
			Msg:  r.Err.Error(),
			Code: r.ErrCode,
		}
	}

	if r.Peers != nil {
		*c.NumberOfPeers = int32(len(r.Peers))
	}

	if r.RawNodeinfo != nil {
		c.RawNodeinfo = r.RawNodeinfo
	}

	n := r.Nodeinfo
	if n != nil {
		*c.SoftwareName = n.SoftwareName()
		*c.SoftwareVersion = n.SoftwareVersion()
		*c.OpenRegistrations = n.IsRegistrationOpen()

		c.TotalUsers = utils.ConvertIntPtrToInt32Ptr(n.TotalUsers())
		c.ActiveHalfyear = utils.ConvertIntPtrToInt32Ptr(n.ActiveUsersHalfyear())
		c.ActiveMonth = utils.ConvertIntPtrToInt32Ptr(n.ActiveUsersMonth())
		c.LocalPosts = utils.ConvertIntPtrToInt32Ptr(n.LocalPosts())
		c.LocalComments = utils.ConvertIntPtrToInt32Ptr(n.LocalComments())
	}

	return c
}

func (r *CrawlResult) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("domain", r.Domain),
		slog.String("error", r.Err.Error()),
		slog.Time("start", r.Start),
		slog.Time("end", r.End),
		slog.String("resolved_ips", fmt.Sprint(r.ResolvedIPs)),
	)
}

func (c *Crawler) Crawl(ctx context.Context, domain string) *CrawlResult {
	slog.InfoContext(ctx, "crawling", "domain", domain)

	r := &CrawlResult{}

	// we capture the end time here
	defer func() {
		r.End = time.Now()
		slog.InfoContext(ctx, "crawled", "domain", domain, "result", r)
	}()

	r.Start = time.Now()
	r.Domain = domain

	// lookup the domain via DNS
	// avoids retrying on such hosts
	ips, err := net.LookupIP(domain)
	r.ResolvedIPs = ips
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			r.Err = err
			r.ErrCode = models.CrawlErrCodeTimeout
		} else {
			r.Err = err
			r.ErrCode = models.CrawlErrCodeDomainNotFound
		}
		return r
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
				r.Err = err
				r.ErrCode = models.CrawlErrCodeTimeout
				return r
			}
		}
		if !canProceed {
			r.Err = errors.New("robots.txt does not allow crawling")
			r.ErrCode = models.CrawlErrCodeBlockedByRobotsTxt
			return r
		} else {
			break
		}
	}
	if err != nil {
		// an error occurred, but we can proceed as if the robots.txt allowed crawling
		slog.ErrorContext(ctx, "failed to acknowledge robots.txt", "error", err)
	}

	// retry for the rest of the requests
	c.client.RetryMax = 3

	nodeInfo, raw, code, err := c.getNodeInfo(ctx, url)
	r.RawNodeinfo = raw
	r.Nodeinfo = nodeInfo
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			r.ErrCode = models.CrawlErrCodeTimeout
		}
		r.Err = err
		r.ErrCode = code
		return r
	}

	peers, code, err := c.GetPeers(ctx, url, nodeInfo)
	r.Peers = peers
	if err != nil {
		if ctx.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			r.ErrCode = models.CrawlErrCodeTimeout
		}
		r.Err = err
		r.ErrCode = code
		return r
	}

	return r
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
		return true, err
	}

	group := robots.FindGroup(c.userAgent)

	// TODO: fill in more endpoints
	canProceed :=
		group.Test(url+"/.well-known/nodeinfo") &&
			group.Test(url+"/api/v1/instance/peers")

	return canProceed, nil
}

// getNodeInfo gets the nodeinfo from the given url.
// returns the nodeinfo and the raw json.
func (c *Crawler) getNodeInfo(ctx context.Context, url string) (nodeinfo.Nodeinfo, []byte, models.CrawlErrCode, error) {
	r, err := retryablehttp.NewRequest("GET", url+"/.well-known/nodeinfo", nil)
	if err != nil {
		return nil, nil, models.CrawlErrCodeInternalError, err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, models.CrawlErrCodeUnreachable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, nil, models.CrawlErrCodeUnreachable, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var w nodeinfo.WellKnown
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return nil, nil, models.CrawlErrCodeInvalidJSON, err
	}

	if len(w.Links) == 0 {
		return nil, nil, models.CrawlErrCodeInvalidNodeinfo, fmt.Errorf("no nodeinfo links found")
	}

	link, nodeInfo, err := nodeinfo.HighestSupported(w)
	if err != nil {
		return nil, nil, models.CrawlErrCodeNodeinfoVersionNotSupportedByCrawl, err
	}

	r, err = retryablehttp.NewRequest("GET", link, nil)
	if err != nil {
		return nil, nil, models.CrawlErrCodeInternalError, err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", c.userAgent)

	resp, err = c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, models.CrawlErrCodeUnreachable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, nil, models.CrawlErrCodeUnreachable, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, b, models.CrawlErrCodeInternalError, err
	}

	err = json.Unmarshal(b, &nodeInfo)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode nodeinfo", "error", err, "body", b)
		return nil, b, models.CrawlErrCodeInvalidNodeinfo, err
	}

	return nodeInfo, b, models.CrawlErrCodeUnknown, nil
}
