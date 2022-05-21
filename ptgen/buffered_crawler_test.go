package ptgen

import (
	"crawler/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetBangumiLinkByNames(t *testing.T) {
	crawler := NewBufferedPTGen(client, kvStorage)
	links, err := crawler.GetBangumiLinkByNames("であいもん", "相合之物", "Deaimon")
	assert.NoError(t, err, "error when getting info")
	for _, v := range links {
		result, err := crawler.GetBangumiDetailByLink(v)
		assert.NoError(t, err, "error when getting link")
		log.Println(GetTextFromDetail(result))
	}
}

var client Client
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
	client, err = NewClient()
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
}
