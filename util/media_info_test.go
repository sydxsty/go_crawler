package util

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetMediaInfo(t *testing.T) {
	info, err := GetMediaInfo("../lib", "../lib", "example.ogg")
	assert.NoError(t, err, "test failure")
	log.Print(info)
}
