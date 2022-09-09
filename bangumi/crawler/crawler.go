package crawler

import (
	"crawler/bangumi"
	"crawler/bangumi/dao"
	"crawler/storage"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

type Bangumi interface {
	GetTagByKeyWord(keyword string) ([]string, error)
	GetAnimeListRawByTag(tag string, page int) ([]interface{}, error)
	GetMiscByTags(ids ...string) []map[string]interface{}
	GetUserNameByTag(ids ...string) []map[string]interface{}
	GetTeamByTag(ids ...string) []map[string]interface{}
}

type BangumiImpl struct {
	client bangumi.Client
	db     storage.KVStorage
}

func NewBangumi(client bangumi.Client, db storage.KVStorage) Bangumi {
	f := &BangumiImpl{
		client: client.Clone(),
		db:     db,
	}
	return f
}

func GetAnimeList(b Bangumi, rawAnimeList []interface{}) ([]*dao.BangumiTorrentInfo, error) {
	var tiList []*dao.BangumiTorrentInfo
	for _, e := range rawAnimeList {
		ti, err := dao.NewTorrentInfoFromMap(e.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		var tags []string
		for _, v := range ti.TagIds {
			tags = append(tags, v.(string))
		}
		ti.InitTorrentDetail(b.GetMiscByTags(tags...))
		tiList = append(tiList, ti)
	}
	return tiList, nil
}

func GetAnimeNameList(b Bangumi, rawAnimeList []interface{}) []map[string]interface{} {
	var tagList []string
	for _, id := range rawAnimeList {
		sl := id.(map[string]interface{})["tag_ids"]
		for _, item := range sl.([]interface{}) {
			tagList = append(tagList, item.(string))
		}
	}
	return b.GetMiscByTags(tagList...)
}

func (b *BangumiImpl) GetTagByKeyWord(keyword string) ([]string, error) {
	resp, err := b.client.SyncPostJson(`api/tag/search`, map[string]interface{}{
		"keywords": true,
		"multi":    true,
		"name":     keyword,
	})
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}
	for _, k := range []string{"success", "found"} {
		if ok1, ok2 := result[k]; !ok2 {
			return nil, errors.New("key-value pair not found, " + k)
		} else {
			if ok3, ok4 := ok1.(bool); !ok3 || !ok4 {
				return nil, errors.New("keyword contains no result, " + k)
			}
		}
	}
	var ids []string
	if tags, ok := result["tag"]; !ok {
		return nil, errors.New("bgm interface is updated, please update your code, " + keyword)
	} else {
		_, ok = tags.([]interface{})
		if !ok {
			return nil, errors.New("bgm interface is updated, please update your code, " + keyword)
		}
		for _, tag := range tags.([]interface{}) {
			_, ok = tag.(map[string]interface{})
			if !ok {
				log.Printf("%s: tag response is not a map", keyword)
				continue
			}
			id, ok := tag.(map[string]interface{})["_id"]
			if !ok {
				log.Printf("%s: tag id not found", keyword)
				continue
			}
			idStr, ok := id.(string)
			if !ok {
				log.Printf("%s: tag id is not a string", keyword)
				continue
			}
			ids = append(ids, idStr)
		}
	}
	return ids, nil
}

func (b *BangumiImpl) GetAnimeListRawByTag(tag string, page int) ([]interface{}, error) {
	if page <= 0 {
		return nil, errors.Errorf("error page index, %d.", page)
	}
	var processFunc func() (*colly.Response, error)
	if tag == "" {
		processFunc = func() (*colly.Response, error) {
			return b.client.SyncVisit(`api/torrent/page/` + strconv.Itoa(page))
		}
	} else {
		processFunc = func() (*colly.Response, error) {
			return b.client.SyncPostJson(`api/torrent/search`, map[string]interface{}{
				"tag_id": []string{tag},
				"p":      page,
			})
		}
	}
	resp, err := processFunc()
	if err != nil {
		return nil, err
	}
	tl := make(map[string]interface{})
	err = json.Unmarshal(resp.Body, &tl)
	if err != nil {
		return nil, err
	}
	return tl["torrents"].([]interface{}), nil
}

func (b *BangumiImpl) GetMiscByTags(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/tag/fetch`)
}

func (b *BangumiImpl) GetUserNameByTag(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/user/fetch`)
}

func (b *BangumiImpl) GetTeamByTag(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/team/fetch`)
}

func (b *BangumiImpl) getBufferedPropertyByTag(ids []string, link string) []map[string]interface{} {
	var resultList []map[string]interface{}
	var animeIDList []string
	for _, id := range ids {
		var r map[string]interface{}
		if err := b.db.Get(id, &r); err != nil {
			// not loaded, get from web
			animeIDList = append(animeIDList, id)
		} else {
			resultList = append(resultList, r)
		}
	}
	// load the rest and save to db
	if len(animeIDList) != 0 {
		tags, err := b.getPropertyByTag(animeIDList, link)
		if err != nil {
			log.Println("error occurred when get from web, ", err)
		}
		for _, r := range tags {
			resultList = append(resultList, r)
			if err = b.db.Put(r[`_id`].(string), r); err != nil {
				log.Println("error occurred when store to db, ", err)
			}
		}
	}
	return resultList
}

func (b *BangumiImpl) getPropertyByTag(id []string, link string) ([]map[string]interface{}, error) {
	resp, err := b.client.SyncPostJson(link, map[string]interface{}{
		"_ids": id,
	})
	if err != nil {
		return nil, err
	}
	var tags []map[string]interface{}
	err = json.Unmarshal(resp.Body, &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
