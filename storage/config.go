package storage

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Username          string
	Password          string
	UseCookie         bool
	CookiePath        string
	QBAddr            string
	QBUsername        string
	QBPassword        string
	LevelDBPath       string
	TorrentPath       string
	ThreadWaterMark   int
	DiscountWaterMark int
}

func LoadConfig(path string) (*Config, error) {
	conf := make(map[interface{}]interface{})
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "config not exist")
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "can not load config")
	}
	cfg := &Config{
		Username:          conf["username"].(string),
		Password:          conf["password"].(string),
		UseCookie:         conf["use_cookie"].(bool),
		CookiePath:        conf["cookie_path"].(string),
		QBAddr:            conf["qb_addr"].(string),
		QBUsername:        conf["qb_username"].(string),
		QBPassword:        conf["qb_password"].(string),
		LevelDBPath:       conf["leveldb_path"].(string),
		TorrentPath:       conf["torrent_path"].(string),
		ThreadWaterMark:   conf["thread_water_mark"].(int),
		DiscountWaterMark: conf["discount_water_mark"].(int),
	}
	if _, err = os.Stat(cfg.TorrentPath); os.IsNotExist(err) {
		_ = os.Mkdir(cfg.TorrentPath, 0666)
	}
	return cfg, nil
}
