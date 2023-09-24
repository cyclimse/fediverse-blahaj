package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyclimse/fediverse-blahaj/internal/models"
	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"
	"github.com/hashicorp/go-retryablehttp"
)

func (c *Crawler) GetPeers(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, models.CrawlErrCode, error) {
	if n == nil {
		return nil, models.CrawlErrCodeInternalError, fmt.Errorf("nodeinfo is nil")
	}
	switch n.SoftwareName() {
	case "mastodon":
		return c.GetPeersMastodon(ctx, url, n)
	}
	return nil, models.CrawlErrCodeSoftwareNotSupportedByCrawler, fmt.Errorf("software not supported by crawler: %s", n.SoftwareName())
}

func (c *Crawler) GetPeersMastodon(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, models.CrawlErrCode, error) {
	r, err := retryablehttp.NewRequest("GET", url+"/api/v1/instance/peers", nil)
	if err != nil {
		return nil, models.CrawlErrCodeInternalError, err
	}

	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, models.CrawlErrCodeUnreachable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, models.CrawlErrCodeUnreachable, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var peers []string
	err = json.NewDecoder(resp.Body).Decode(&peers)
	if err != nil {
		return nil, models.CrawlErrCodeInvalidJSON, err
	}

	return peers, models.CrawlErrCodeUnknown, nil
}
