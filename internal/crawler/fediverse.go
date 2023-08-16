package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"
	"github.com/hashicorp/go-retryablehttp"
)

func (c *Crawler) GetPeers(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, CrawlError) {
	if n == nil {
		return nil, newCrawlInternalError(fmt.Errorf("nodeinfo is nil"))
	}
	switch n.GetSoftwareName() {
	case "mastodon":
		return c.GetPeersMastodon(ctx, url, n)
	}
	return nil, errUnsupportedSoftware
}

func (c *Crawler) GetPeersMastodon(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, CrawlError) {
	r, err := retryablehttp.NewRequest("GET", url+"/api/v1/instance/peers", nil)
	if err != nil {
		return nil, newCrawlInternalError(err)
	}

	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, errNetworkError.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errNetworkError.Wrap(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	var peers []string
	err = json.NewDecoder(resp.Body).Decode(&peers)
	if err != nil {
		return nil, newCrawlInternalError(err)
	}

	return peers, nil
}
