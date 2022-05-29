package ptgen

import (
	"crawler/storage"
	"github.com/pkg/errors"
	"log"
)

type BufferedPTGenImpl struct {
	pg PTGen
	db storage.KVStorage
}

func NewBufferedPTGen(client Client, db storage.KVStorage) PTGen {
	return &BufferedPTGenImpl{
		pg: NewPTGen(client),
		db: db,
	}
}

func (b BufferedPTGenImpl) GetBangumiLinkByNames(jpnName string, names ...string) ([]*BangumiLinkDetail, error) {
	for _, name := range append([]string{jpnName}, names...) {
		result, err := b.GetBangumiLinkByName(name)
		if err == nil {
			return result, nil
		}
		log.Println("this name contains no result, switch to another name, ", name, err)
	}
	return nil, errors.New("can not get any valid result")
}

func (b BufferedPTGenImpl) GetBangumiLinkByName(name string) ([]*BangumiLinkDetail, error) {
	if name == "" {
		return nil, errors.New("query name is empty")
	}
	results, err := loadPTGenResult(b.db, name)
	if err == nil {
		return results, nil
	}
	results, err = b.pg.GetBangumiLinkByName(name)
	if err != nil {
		return nil, err
	}
	err = savePTGenResult(b.db, name, results)
	if err != nil {
		log.Println("pt-gen result save failed, ", err)
	}
	return results, nil
}

func (b BufferedPTGenImpl) GetBangumiInfoByLink(link string) (map[string]interface{}, error) {
	results, err := loadPTGenLink(b.db, link)
	if err == nil {
		success, ok := results["success"].(bool)
		if ok && success == true {
			return results, nil
		}
	}
	results, err = b.pg.GetBangumiInfoByLink(link)
	if err != nil {
		return nil, err
	}
	err = savePTGenLink(b.db, link, results)
	if err != nil {
		log.Println("pt-gen link result save failed, ", err)
	}
	return results, nil
}

func savePTGenResult(db storage.KVStorage, chsName string, result []*BangumiLinkDetail) error {
	err := db.Put(`pt_gen_`+chsName, result)
	if err != nil {
		return err
	}
	return nil
}

func loadPTGenResult(db storage.KVStorage, chsName string) ([]*BangumiLinkDetail, error) {
	var result []*BangumiLinkDetail
	err := db.Get(`pt_gen_`+chsName, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func savePTGenLink(db storage.KVStorage, link string, result map[string]interface{}) error {
	err := db.Put(`pt_gen_`+link, result)
	if err != nil {
		return err
	}
	return nil
}

func loadPTGenLink(db storage.KVStorage, link string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := db.Get(`pt_gen_`+link, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
