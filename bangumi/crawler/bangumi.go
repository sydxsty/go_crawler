package crawler

import (
	"crawler/bangumi/dao"
	"crawler/qbt"
	"log"
	"os"
)

// ScanBangumiTorrent call the callback when get a torrent info from website
func ScanBangumiTorrent(bgm Bangumi, postTorrentFunc func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	ral, err := bgm.GetLatestAnimeListRaw()
	if err != nil {
		return err
	}
	al, err := GetAnimeList(bgm, ral)
	if err != nil {
		return err
	}
	for _, anime := range al {
		postTorrentFunc(anime)
	}
	log.Println("all torrent finished scanning, return")
	return nil
}

func DownloadBangumiTorrent(url string, d Downloader, q qbt.WEBUIHelper) ([]byte, error) {
	data, err := d.DownloadTorrentFromUrl(url)
	if err != nil {
		return nil, err
	}
	err = q.AddTorrentFromData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DownloadBangumiTorrentToFile(url string, path string, hash string, d Downloader, q qbt.WEBUIHelper) ([]byte, error) {
	data, err := d.DownloadTorrentFromUrl(url)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path+hash+`.torrent`, data, 0666)
	if err != nil {
		return nil, err
	}
	err = q.AddTorrentFromFile(path, hash+`.torrent`)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func LoadTorrentFromFile(path string, hash string) ([]byte, error) {
	data, err := os.ReadFile(path + hash + `.torrent`)
	if err != nil {
		return nil, err
	}
	return data, nil
}
