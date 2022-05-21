package crawler

import (
	"crawler/bangumi"
	"crawler/bangumi/dao"
	"crawler/storage"
	"encoding/json"
	"log"
)

type Bangumi interface {
	GetLatestAnimeListRaw() ([]interface{}, error)
	GetRecentAnimeListRaw() ([]interface{}, error)
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

func (b *BangumiImpl) GetLatestAnimeListRaw() ([]interface{}, error) {
	resp, err := b.client.SyncVisit(`api/torrent/latest`)
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

func (b *BangumiImpl) GetRecentAnimeListRaw() ([]interface{}, error) {
	resp, err := b.client.SyncVisit(`api/torrent/recent`)
	if err != nil {
		return nil, err
	}
	var tl []interface{}
	err = json.Unmarshal(resp.Body, &tl)
	if err != nil {
		return nil, err
	}
	return tl, nil
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

func (b *BangumiImpl) getBufferedPropertyByTag(ids []string, url string) []map[string]interface{} {
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
		tags, err := b.getPropertyByTag(animeIDList, url)
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

func (b *BangumiImpl) getPropertyByTag(id []string, url string) ([]map[string]interface{}, error) {
	postData := make(map[string]interface{})
	postData["_ids"] = id
	raw, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	resp, err := b.client.SyncPostRaw(url, raw)
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
