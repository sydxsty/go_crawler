package crawler

import (
	"crawler/bangumi/dao"
	"crawler/qbt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

// ScanBangumiTorrent call the callback when get a torrent info from website
func ScanBangumiTorrent(bgm Bangumi, postTorrentFunc func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	for i := 1; i < 5; i++ {
		log.Printf("scanning page %d", i)
		tl, err := bgm.GetAnimeListRawByTag("", i)
		if err != nil {
			return err
		}
		al, err := GetAnimeList(bgm, tl)
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

func CrawlAllTorrents(bgm Bangumi, keyword string, postTorrentFunc func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	var tag string
	{
		tags, err := bgm.GetTagByKeyWord(keyword)
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			return errors.New("no valid search result")
		}
		tag = tags[0]
	}
	for i := 1; ; i++ {
		log.Printf("scanning page %d", i)
		tl, err := bgm.GetAnimeListRawByTag(tag, i)
		if err != nil {
			return err
		}
		if len(tl) == 0 {
			break
		}
		al, err := GetAnimeList(bgm, tl)
		if err != nil {
			return err
		}
		for _, anime := range al {
			postTorrentFunc(anime)
		}
		log.Println("wait 10 sec to continue")
		time.Sleep(time.Second * 10)
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
