package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"
	"github.com/hashicorp/go-retryablehttp"
)

func (c *Crawler) GetPeers(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, error) {
	if n == nil {
		return nil, fmt.Errorf("no software information")
	}
	switch n.GetSoftwareName() {
	case "mastodon":
		return c.GetPeersMastodon(ctx, url, n)
	}
	return nil, fmt.Errorf("unknown software: %s", n.GetSoftwareName())
}

func (c *Crawler) GetPeersMastodon(ctx context.Context, url string, n nodeinfo.Nodeinfo) ([]string, error) {
	r, err := retryablehttp.NewRequest("GET", url+"/api/v1/instance/peers", nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get peers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get peers, status code: %d", resp.StatusCode)
	}

	var peers []string
	err = json.NewDecoder(resp.Body).Decode(&peers)
	if err != nil {
		return nil, fmt.Errorf("failed to decode peers: %w", err)
	}

	return peers, nil
}
