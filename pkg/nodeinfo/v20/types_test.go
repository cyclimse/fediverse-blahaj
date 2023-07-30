package v20

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadNodeInfoFixture(t *testing.T, path string) []byte {
	f, err := os.Open(path)
	require.NoError(t, err)

	b, err := io.ReadAll(f)
	require.NoError(t, err)

	return b
}

func TestNodeInfoValidateOnExampleNodeInfo(t *testing.T) {
	b := loadNodeInfoFixture(t, "./fixtures/example-nodeinfo.json")

	var nodeInfo Nodeinfo
	err := json.Unmarshal(b, &nodeInfo)
	require.NoError(t, err)
}

func TestNodeInfoValidateOnInvalidExampleNodeInfo(t *testing.T) {
	// Is missing the required field "version"
	b := loadNodeInfoFixture(t, "./fixtures/example-nodeinfo-invalid.json")

	var nodeInfo Nodeinfo
	err := json.Unmarshal(b, &nodeInfo)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "field version in NodeinfoSchemaJsonServer: required")
}
