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

func (t *TorrentManager) ResetTorrentPostedState(animeInfoHash string) bool {
	err := t.setState(animeInfoHash, "posted", false)
	if err != nil {
		return false
	}
	return true
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
	teams := info.MustGetTeam()
	if teams == nil {
		return errors.New("no team name found in info")
	}
	name := info.MustGetCHSName()
	if len(name) == 0 {
		return errors.New("torrent name is empty")
	}
	episode := info.MustGetEpisode()
	if v := regexp.MustCompile(`(720)`).FindAllString(info.Title, -1); len(v) != 0 {
		return errors.New("720p, skip download, " + v[0])
	}

	if v := regexp.MustCompile(`((简体)|(简日)|(简繁))`).FindAllString(info.Title, -1); len(v) == 0 {
		// TODO: have bugs, consider using machine learning
		if v := regexp.MustCompile(`((繁體)|(繁日)|(CHT)|(BIG5)|(繁体))`).FindAllString(info.Title, -1); len(v) != 0 {
			return errors.New("is not sc, skip download, " + v[0])
		}
	}
	detail := make(map[string]interface{})
	// get the uploaded team of current episode
	v, err := t.getState(name, episode)
	if err == nil {
		if tmp, ok := v.(map[string]interface{}); ok && tmp != nil {
			detail = tmp
		}
	}
	for _, teamName := range teams {
		if _, ok := detail[teamName]; ok {
			return errors.New("we have already downloaded the same torrent")
		}
		detail[teamName] = true
	}
	err = t.setState(name, episode, detail)
	if err != nil {
		return err
	}
	return nil
}
