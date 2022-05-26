package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetThreadDetail get all torrent pages from a thread
func TestGetThreadDetail(t *testing.T) {
	detail := NewThreadDetail(client)
	allFloor, err := detail.GetFloorDetailFromThread(`thread-1698281-1-1.html`)
	assert.NoError(t, err, "error get allFloor in a thread")
	assert.True(t, len(allFloor) >= 1, "floor size < 1, impossible")
	FilteredFloor, err := detail.GetFloorDetailFromThread(`thread-1698281-1-1.html`)
	assert.NoError(t, err, "error FilteredFloor")
	assert.True(t, len(FilteredFloor) == 1, "torrent link > 1, impossible")
}
