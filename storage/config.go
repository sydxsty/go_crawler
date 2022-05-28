package storage

import (
	"gopkg.in/yaml.v3"
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
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		if err := yaml.NewDecoder(f).Decode(conf); err != nil {
			return nil, err
		}
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
	return cfg, nil
}
