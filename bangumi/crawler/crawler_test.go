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
	processTorrentListFunc := func(torrentList []interface{}) {
		al, err := GetAnimeList(crawler, torrentList)
		assert.NoError(t, err, "we must process the torrent list")
		assert.True(t, len(al) > 0, "we must process at least one torrent successfully")
	}
	// test get latest
	tl, err := crawler.GetAnimeListRawByTag(nil, 1)
	assert.NoError(t, err, "Call the api to return a list of torrent")
	processTorrentListFunc(tl)
	// test get list with keyword
	tags, err := crawler.GetTagByKeyWord("lycoris")
	assert.NoError(t, err, "can not get tag for a torrent")
	assert.True(t, len(tags) > 0, "result must contains at least one tag")
	tl, err = crawler.GetAnimeListRawByTag(tags[:1], 1)
	assert.NoError(t, err, "Call the api to return at least one tag")
	processTorrentListFunc(tl)
}

func TestGetAnimeNameList(t *testing.T) {
	crawler := NewBangumi(client, kvStorage)
	tl, err := crawler.GetAnimeListRawByTag(nil, 1)
	assert.NoError(t, err, "error GetLatestAnimeListRaw")
	bl := GetAnimeNameList(crawler, tl)
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
