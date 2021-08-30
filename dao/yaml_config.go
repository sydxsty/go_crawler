package dao

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v2"
	"log"
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

var YAMLConfig *Config
var dbHandle *leveldb.DB

func init() {
	err := ReadYamlConfig("./data/config.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}
	dbHandle, err = leveldb.OpenFile(YAMLConfig.LevelDBPath, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ReadYamlConfig(path string) error {
	conf := make(map[interface{}]interface{})
	if f, err := os.Open(path); err != nil {
		return err
	} else {
		if err := yaml.NewDecoder(f).Decode(conf); err != nil {
			return err
		}
	}
	YAMLConfig = &Config{
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
	return nil
}

func LoadFromDB(uid string, data interface{}) error {
	value, err := dbHandle.Get([]byte(uid), nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(value, data); err != nil {
		return err
	}
	return nil
}

func SaveToDB(uid string, data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := dbHandle.Put([]byte(uid), raw, nil); err != nil {
		return err
	}
	return nil
}
