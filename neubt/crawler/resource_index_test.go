package crawler

import (
	"crawler/neubt/dao"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetResourceIndex get all torrent pages from resource index
func TestGetResourceIndex(t *testing.T) {
	ri := NewResourceIndex(client)
	list, err := ri.GetResourceIndex()
	assert.NoError(t, err, "error phrase resource index")
	result, err := dao.NodeListToTorrentInfoList(kvStorage, list)
	assert.NoError(t, err, "error get torrent info")
	assert.True(t, len(result) > 0)
}
