package crawler_test

import (
	"context"
	"testing"

	"github.com/cyclimse/fediverse-blahaj/internal/crawler"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v20"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPeersMastodon(t *testing.T) {
	c := crawler.New("test")
	nodeInfo := nodeinfo.Nodeinfo{
		Software: nodeinfo.NodeinfoSoftware{
			Name: "mastodon",
		},
	}

	peers, code, err := c.GetPeersMastodon(context.Background(), "https://mastodon.social", &nodeInfo)
	require.NoError(t, err)
	assert.Equal(t, code, models.CrawlErrCodeUnknown)
	assert.Greater(t, len(peers), 100)
}
