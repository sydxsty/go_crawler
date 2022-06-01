package ptgen

import (
	"crawler/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetBangumiLinkByNames(t *testing.T) {
	crawler, err := NewBufferedPTGen(kvStorage)
	assert.NoError(t, err, "init failure")
	links, err := crawler.GetBangumiLinkByNames("であいもん", "相合之物", "Deaimon")
	assert.NoError(t, err, "error when getting info")
	for _, v := range links {
		result, err := crawler.GetBangumiInfoByLink(v.Link)
		assert.NoError(t, err, "error when getting link")
		log.Println(GetDetailFromInfo(result))
	}
}

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
}
