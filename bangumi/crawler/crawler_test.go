package crawler

import (
	"crawler/bangumi"
	"crawler/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetAnimeList(t *testing.T) {
	crawler := NewBangumi(client, kvStorage)
	tl, err := crawler.GetLatestAnimeListRaw()
	assert.NoError(t, err, "error GetLatestAnimeListRaw")
	bl, err := GetAnimeList(crawler, tl)
	assert.NoError(t, err, "error GetLatestAnimeList")
	assert.True(t, len(bl) > 0)
	tl, err = crawler.GetRecentAnimeListRaw()
	assert.NoError(t, err, "error GetRecentAnimeListRaw")
	bl, err = GetAnimeList(crawler, tl)
	assert.NoError(t, err, "error GetLatestAnimeList")
	assert.True(t, len(bl) > 0)
}

func TestGetAnimeNameList(t *testing.T) {
	crawler := NewBangumi(client, kvStorage)
	tl, err := crawler.GetLatestAnimeListRaw()
	assert.NoError(t, err, "error GetLatestAnimeListRaw")
	bl := GetAnimeNameList(crawler, tl)
	assert.True(t, len(bl) > 0)
	tl, err = crawler.GetRecentAnimeListRaw()
	assert.NoError(t, err, "error GetRecentAnimeListRaw")
	bl = GetAnimeNameList(crawler, tl)
	assert.True(t, len(bl) > 0)
}

var client bangumi.Client
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
	client, err = bangumi.NewClient()
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
}
