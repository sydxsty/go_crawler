package qbt

import (
	"crawler/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestWEBUI(t *testing.T) {
	infoHash := `6ce36d8b6688fd1439ba0a9427c6ddcdd85bdcde`
	config, err := storage.LoadConfig("./data/config.yaml")
	assert.NoError(t, err, "can not load config")
	webui, err := NewWEBUIHelper(config.QBAddr, config.QBUsername, config.QBPassword)
	assert.NoError(t, err, "can not connect to webui")
	assert.True(t, webui.Contains(infoHash), "please enter the hash")
	torrent, files, err := webui.GetTorrentDetail(infoHash)
	assert.NoError(t, err, "can not find the torrent files of specific hash")
	log.Print(torrent)
	log.Print(files)
}
