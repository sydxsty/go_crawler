package main

import (
	"crawler/bangumi/crawler"
	"crawler/bangumi/dao"
	"log"
	"strings"
	"time"
)

func main() {
	p := NewPoster()
	for {
		// KeyWordCrawler(p)
		DefaultIndexCrawler(p)
		log.Println("wait 600 sec to recheck")
		time.Sleep(time.Second * 600)
	}
}

func DefaultIndexCrawler(p *Poster) {
	err := crawler.ScanBangumiTorrent(p.bgm, p.BGMSearchCallback)
	if err != nil {
		log.Println("can not load bangumi latest torrents")
		time.Sleep(time.Second * 60)
	}
}

func KeyWordCrawler(p *Poster) {
	priorityTeamList := []string{"喵萌奶茶屋", "LoliHouse"}
	animeList := []string{"杜鹃的婚约"}
	keywordInTitle := ""
	if len(animeList) == 0 {
		return
	}
	callbackWrapper := func(ti *dao.BangumiTorrentInfo) {
		if strings.Count(ti.Title, keywordInTitle) == 0 {
			log.Printf("skip this torrent %s, keyword %s.", ti.Title, keywordInTitle)
			return
		}
		p.BGMSearchCallback(ti)
	}
	for _, aniName := range animeList {
		for _, teamName := range priorityTeamList {
			err := crawler.CrawlAllTorrents(p.bgm, []string{teamName, aniName}, callbackWrapper)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("target torrent found")
			break
		}
	}
}
