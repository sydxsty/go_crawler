package dao

import (
	"crawler/storage"
	"errors"
	"regexp"
)

type TorrentManager struct {
	db storage.KVStorage
}

func NewTorrentManager(db storage.KVStorage) *TorrentManager {
	return &TorrentManager{
		db: db,
	}
}

func (t *TorrentManager) getState(animeInfoHash string, str string) (interface{}, error) {
	var state map[string]interface{}
	err := t.db.Get("bgm_anime_state_"+animeInfoHash, &state)
	if err != nil {
		return nil, err
	}
	if v, ok := state[str]; !ok {
		return nil, errors.New("state not exist")
	} else {
		return v, nil
	}
}

func (t *TorrentManager) setState(animeInfoHash string, str string, val interface{}) error {
	state := make(map[string]interface{})
	_ = t.db.Get("bgm_anime_state_"+animeInfoHash, &state)
	state[str] = val
	err := t.db.Put("bgm_anime_state_"+animeInfoHash, state)
	if err != nil {
		return err
	}
	return nil
}

func (t *TorrentManager) TorrentIsPosted(animeInfoHash string) bool {
	state, err := t.getState(animeInfoHash, "posted")
	if err != nil {
		return false
	}
	v, ok := state.(bool)
	if !ok {
		return false
	}
	return v
}

func (t *TorrentManager) SetTorrentPostedState(animeInfoHash string) bool {
	err := t.setState(animeInfoHash, "posted", true)
	if err != nil {
		return false
	}
	return true
}

// CanDownloadFromBangumi this is a filter to avoid download same kind of torrents from site
func (t *TorrentManager) CanDownloadFromBangumi(info *BangumiTorrentInfo) error {
	name, err := info.GetCHNName()
	if err != nil {
		return errors.New("torrent name is empty")
	}
	if v := regexp.MustCompile(`(720)`).FindAllString(info.Title, -1); len(v) != 0 {
		return errors.New("720p, skip download, " + v[0])
	}
	// TODO: have bugs, consider using machine learning
	if v := regexp.MustCompile(`((繁體)|(繁日)|(CHT)|(BIG5))`).FindAllString(info.Title, -1); len(v) != 0 {
		return errors.New("is not sc, skip download, " + v[0])
	}

	if v := regexp.MustCompile(`((NC-Raws)|(NaN-Raws)|(ANi)|(TD-RAWS)|(国漫))`).FindAllString(info.Title, -1); len(v) != 0 {
		return errors.New("wrong team, skip download, " + v[0])
	}

	if v := regexp.MustCompile(`((喵萌奶茶屋)|(LoliHouse)|(喵萌Production)|(字幕组))`).FindAllString(info.Title, -1); len(v) == 0 {
		return errors.New("does not get correct team" + info.Title)
	}

	detail := make(map[string]interface{})
	// get the uploaded team of current episode
	v, err := t.getState(name, info.Detail.Episode)
	if err == nil {
		if tmp, ok := v.(map[string]interface{}); ok && tmp != nil {
			detail = tmp
		}
	}

	if _, ok := detail[info.Detail.TeamName]; ok {
		return errors.New("we have already downloaded the same torrent")
	}
	detail[info.Detail.TeamName] = true
	err = t.setState(name, info.Detail.Episode, detail)
	if err != nil {
		return err
	}
	return nil
}
