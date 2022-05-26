package ptgen

import (
	"crawler/storage"
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

func (b BufferedPTGenImpl) GetBangumiLinkByNames(jpnName string, names ...string) (map[string]string, error) {
	results, err := loadPTGenResult(b.db, jpnName)
	if err == nil {
		return results, nil
	}
	results, err = b.pg.GetBangumiLinkByNames(jpnName, names...)
	if err != nil {
		return nil, err
	}
	err = savePTGenResult(b.db, jpnName, results)
	if err != nil {
		log.Println("pt-gen result save failed, ", err)
	}
	return results, nil
}

func (b BufferedPTGenImpl) GetBangumiLinkByName(name string) (map[string]string, error) {
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

func (b BufferedPTGenImpl) GetBangumiDetailByLink(link string) (map[string]interface{}, error) {
	results, err := loadPTGenLink(b.db, link)
	if err == nil {
		success, ok := results["success"].(bool)
		if ok && success == true {
			return results, nil
		}
	}
	results, err = b.pg.GetBangumiDetailByLink(link)
	if err != nil {
		return nil, err
	}
	err = savePTGenLink(b.db, link, results)
	if err != nil {
		log.Println("pt-gen link result save failed, ", err)
	}
	return results, nil
}

func savePTGenResult(db storage.KVStorage, chsName string, result map[string]string) error {
	err := db.Put(`pt_gen_`+chsName, result)
	if err != nil {
		return err
	}
	return nil
}

func loadPTGenResult(db storage.KVStorage, chsName string) (map[string]string, error) {
	var result map[string]string
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
