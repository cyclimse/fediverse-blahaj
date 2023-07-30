package crawler_test

import (
	"context"
	"testing"

	"github.com/cyclimse/fediverse-blahaj/internal/crawler"
	nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v20"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPeersMastodon(t *testing.T) {
	c := crawler.New()
	nodeInfo := nodeinfo.Nodeinfo{
		Software: nodeinfo.NodeinfoSoftware{
			Name: "mastodon",
		},
	}

	peers, err := c.GetPeersMastodon(context.Background(), "https://mastodon.social", &nodeInfo)
	require.NoError(t, err)
	assert.Greater(t, len(peers), 100)
}
