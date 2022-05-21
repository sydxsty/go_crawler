package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetForum get all torrent pages from a forum
func TestGetForumDetail(t *testing.T) {
	detail := NewForumDetail(client)
	allFloor, err := detail.GetFloorDetailFromForum(`thread-1698281-1-1.html`)
	assert.NoError(t, err, "error get allFloor in a forum")
	assert.True(t, len(allFloor) >= 1, "floor size < 1, impossible")
	FilteredFloor, err := detail.GetFloorDetailFromForum(`thread-1698281-1-1.html`)
	assert.NoError(t, err, "error FilteredFloor")
	assert.True(t, len(FilteredFloor) == 1, "torrent link > 1, impossible")
}
