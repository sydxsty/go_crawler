package dao

import "crawler/storage"

type TorrentManager struct {
	db storage.KVStorage
}

func NewTorrentManager(db storage.KVStorage) *TorrentManager {
	return &TorrentManager{
		db: db,
	}
}

func (t *TorrentManager) getState(animeInfoHash string, str string) bool {
	var state map[string]bool
	err := t.db.Get("bgm_anime_state_"+animeInfoHash, &state)
	if err != nil {
		return false
	}
	if v, ok := state[str]; !ok {
		return false
	} else {
		return v
	}
}

func (t *TorrentManager) setState(animeInfoHash string, str string, val bool) bool {
	state := make(map[string]bool)
	t.db.Get("bgm_anime_state_"+animeInfoHash, &state)
	state[str] = val
	err := t.db.Put("bgm_anime_state_"+animeInfoHash, state)
	if err != nil {
		return false
	}
	return true
}

func (t *TorrentManager) TorrentIsPosted(animeInfoHash string) bool {
	return t.getState(animeInfoHash, "posted")
}

func (t *TorrentManager) SetTorrentPostedState(animeInfoHash string) bool {
	return t.setState(animeInfoHash, "posted", true)
}
