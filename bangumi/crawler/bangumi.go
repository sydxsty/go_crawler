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
func ScanBangumiTorrent(bgm Bangumi, callback func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	for i := 1; i < 5; i++ {
		log.Printf("scanning page %d", i)
		tl, err := bgm.GetAnimeListRawByTag(nil, i)
		if err != nil {
			return err
		}
		al, err := GetAnimeList(bgm, tl)
		if err != nil {
			return err
		}
		for _, anime := range al {
			callback(anime)
		}
		log.Println("wait 300 sec to continue")
		time.Sleep(time.Second * 300)
	}
	log.Println("all torrent finished scanning, return")
	return nil
}

func CrawlAllTorrents(bgm Bangumi, keywords []string, callback func(*dao.BangumiTorrentInfo)) error {
	// init crawler
	var tags []string
	for _, keyword := range keywords {
		res, err := bgm.GetTagByKeyWord(keyword)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			log.Printf("no valid search result for keyword: %s", keyword)
			continue
		}
		tags = append(tags, res[0])
	}
	for i := 1; ; i++ {
		log.Printf("scanning page %d", i)
		tl, err := bgm.GetAnimeListRawByTag(tags, i)
		if err != nil {
			return err
		}
		if len(tl) == 0 {
			if i == 1 {
				return errors.New("no valid result found")
			}
			break
		}
		al, err := GetAnimeList(bgm, tl)
		if err != nil {
			return err
		}
		for _, anime := range al {
			callback(anime)
		}
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
