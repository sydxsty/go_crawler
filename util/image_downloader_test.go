package util

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDownloadImage(t *testing.T) {
	id, err := NewImageDownloader("https://lain.bgm.tv/pic/cover/l/de/4a/329906_hmtVD.jpg")
	assert.NoError(t, err, "error init ImageDownloader")
	data, fileType, err := id.Download()
	assert.NoError(t, err, "error download image")
	assert.True(t, fileType == "jpg")
	assert.True(t, len(data) > 1000000)
	log.Print(data, fileType)
}
