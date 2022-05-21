package crawler

import (
	"crawler/neubt"
	"crawler/neubt/dao"
	"crawler/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// TestGetForumList get all forum from home page
func TestGetForumList(t *testing.T) {
	forum := NewForum(client)
	list, err := forum.GetForumList()
	assert.NoError(t, err, "error phrase home page")
	_, ok := list["高清电影"]
	assert.True(t, ok, "高清电影 not found")
	_, ok = list["卡通动漫"]
	assert.True(t, ok, "卡通动漫 not found")
}

// TestGetForum get all torrent pages from a forum
func TestGetForum(t *testing.T) {
	forum := NewForum(client)
	list, err := forum.GetForum(`forum-13-1.html`)
	assert.NoError(t, err, "error phrase home page")
	result, err := dao.NodeListToTorrentInfoList(kvStorage, list)
	assert.NoError(t, err, "error get torrent info")
	assert.True(t, len(result) > 0)
}

var client neubt.Client
var kvStorage storage.KVStorage

func init() {
	cfg, err := storage.LoadConfig("./data/config.yaml")
	if err != nil {
		log.Fatal("error load config.")
	}
	kvStorage, err = storage.NewKVStorage(cfg.LevelDBPath)
	if err != nil {
		log.Fatal(err, "error load db.")
	}
	client, err = neubt.NewClient(kvStorage)
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
	err = client.LoadCookie(cfg.CookiePath)
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
}
