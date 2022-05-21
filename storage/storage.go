package storage

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

type KVStorage interface {
	GetRaw(key string) ([]byte, error)
	PutRaw(key string, value []byte) error
	Get(key string, value interface{}) error
	Put(key string, value interface{}) error
}

type KVStorageImpl struct {
	db *leveldb.DB
}

func NewKVStorage(dbPath string) (KVStorage, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	kvs := &KVStorageImpl{
		db: db,
	}
	return kvs, nil
}

func (k *KVStorageImpl) GetRaw(key string) ([]byte, error) {
	return k.db.Get([]byte(key), nil)
}

func (k *KVStorageImpl) PutRaw(key string, value []byte) error {
	return k.db.Put([]byte(key), value, nil)
}

func (k *KVStorageImpl) Get(key string, value interface{}) error {
	rawValue, err := k.GetRaw(key)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rawValue, value); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (k *KVStorageImpl) Put(key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return k.PutRaw(key, raw)
}
