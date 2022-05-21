package neubt

import (
	"crawler/qbt"
	"crawler/storage"
	"log"
)

type NeuBT struct {
	Config *storage.Config
	KVS    storage.KVStorage
	Client Client
	Webui  qbt.WEBUIHelper
}

func NewNeuBT() *NeuBT {
	n := &NeuBT{}
	var err error
	n.Config, err = storage.LoadConfig("./data/config.yaml")
	if err != nil {
		log.Fatal("error load config.")
	}
	n.KVS, err = storage.NewKVStorage(n.Config.LevelDBPath)
	if err != nil {
		log.Fatal(err, "error load db.")
	}
	n.Client, err = NewClient(n.KVS)
	if err != nil {
		log.Fatal(err, "can not init client before startup")
	}
	log.Println("try to use cookie login")
	err = n.Client.LoadCookie(n.Config.CookiePath)
	if err != nil {
		log.Println(err, "cookie not found, try password login")
		err = n.Client.Login(n.Config.Username, n.Config.Password)
		if err != nil {
			log.Fatal(err, "password login failed")
		}
		err = n.Client.SaveCookie(n.Config.CookiePath)
		if err != nil {
			log.Fatal(err, "save cookie failed")
		}
	}
	log.Println("login successfully")
	n.Webui, err = qbt.NewWEBUIHelper(n.Config.QBAddr, n.Config.QBUsername, n.Config.QBPassword)
	return n
}
