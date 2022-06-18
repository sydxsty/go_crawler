package crawler

import (
	"crawler/bangumi/dao"
	"crawler/qbt"
	"log"
	"os"
	"time"
)

// ScanBangumiTorrent call the callback when get a torrent info from website
func ScanBangumiTorrent(bgm Bangumi, postTorrentFunc func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	for i := 1; i < 5; i++ {
		log.Printf("scanning page %d", i)
		ral, err := bgm.GetAnimeListRaw(i)
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
		log.Println("wait 600 sec to continue")
		time.Sleep(time.Second * 600)
	}
	log.Println("all torrent finished scanning, return")
	return nil
}

func DownloadBangumiTorrent(link string, d Downloader, q qbt.WEBUIHelper) ([]byte, error) {
	data, err := d.DownloadTorrentFromUrl(link)
	if err != nil {
		return nil, err
	}
	err = q.AddTorrentFromData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DownloadBangumiTorrentToFile(link string, path string, hash string, d Downloader, q qbt.WEBUIHelper) ([]byte, error) {
	data, err := d.DownloadTorrentFromUrl(link)
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
