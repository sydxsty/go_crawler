package bgmtv

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSearchBangumi(t *testing.T) {
	result, err := SearchBangumi("BIRDIE WING -Golf Girls’ Story")
	assert.NoError(t, err, "search failure")
	log.Println(result)
}

func TestGetBangumi(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err, "init failure")
	result, err := GenBangumi(client, "https://bgm.tv/subject/329906")
	assert.NoError(t, err, "search failure")
	log.Println(result)
}
