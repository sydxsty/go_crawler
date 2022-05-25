package neubt

import (
	"crawler/storage"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

var client Client
var cfg *storage.Config

func TestLoginWithPassword(t *testing.T) {
	var resp *colly.Response
	client.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	err := client.Login(os.Getenv("username"), os.Getenv("password"))
	assert.NoError(t, err, "login failure")
	assert.True(t, strings.Index(string(resp.Body), `欢迎您回来`) != -1, "login unsuccessful")
	err = client.SaveCookie(cfg.CookiePath)
	assert.NoError(t, err, "save cookie failure")
}

func TestLoginWithCookie(t *testing.T) {
	err := client.LoadCookie(cfg.CookiePath)
	assert.NoError(t, err, "load cookie failure")
	resp, err := client.SyncVisit(`plugin.php?id=neubt_resourceindex`)
	assert.NoError(t, err, "visit index page failure")
	assert.False(t, strings.Index(string(resp.Body), `登录`) != -1, "cookie login unsuccessful")
}

func TestClone(t *testing.T) {
	err := client.LoadCookie(cfg.CookiePath)
	assert.NoError(t, err, "load cookie failure")
	client2 := client.Clone()
	resp, err := client2.SyncVisit(`plugin.php?id=neubt_resourceindex`)
	assert.NoError(t, err, "visit index page failure")
	assert.False(t, strings.Index(string(resp.Body), `登录`) != -1, "cookie login unsuccessful")
}

func init() {
	var err error
	cfg, err = storage.LoadConfig("./data/config.yaml")
	if err != nil {
		log.Fatal("error load config.")
	}
	kvStorage, err := storage.NewKVStorage(cfg.LevelDBPath)
	if err != nil {
		log.Fatal(err, "error load db.")
	}
	client, err = NewClient(kvStorage)
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
}
