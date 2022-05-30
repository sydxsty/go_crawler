package util

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetMediaInfo(t *testing.T) {
	info, err := GetMediaInfo("./lib", "./lib", "example.ogg")
	assert.NoError(t, err, "test failure")
	log.Print(info)
}

func TestGetMediaImage(t *testing.T) {
	bytes, err := GetMediaImage("./lib", "D:\\迅雷下载", "[Comicat&KissSub][SPY×FAMILY][07][1080P][GB][MP4].mp4")
	assert.NoError(t, err, "test failure")
	log.Print(len(bytes))
}
