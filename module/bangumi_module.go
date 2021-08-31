package module

import (
	"bytes"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"io"
	"log"
	"net/url"
	"os"
)

type Bangumi struct {
	scraperModuleImpl
}

func NewBangumiModule() *Bangumi {
	m := &Bangumi{}
	m.init()
	m.domain, _ = url.Parse("https://bangumi.moe")
	return m
}

func (b *Bangumi) getClonedCollector() *colly.Collector {
	clone := b.collector.Clone()
	clone.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "bangumi.moe")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`)
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Origin", `https://bangumi.moe`)
		r.Headers.Set("Referer", `https://bangumi.moe/`)
		r.Headers.Set("Accept-Encoding", "deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
	return clone
}

func (b *Bangumi) GetAnimeMiscByTag(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/tag/fetch`)
}

func (b *Bangumi) GetUserNameByTag(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/user/fetch`)
}

func (b *Bangumi) GetTeamByTag(ids ...string) []map[string]interface{} {
	return b.getBufferedPropertyByTag(ids, `api/team/fetch`)
}

func (b *Bangumi) getBufferedPropertyByTag(ids []string, url string) []map[string]interface{} {
	var resultList []map[string]interface{}
	var animeIDList []string
	for _, id := range ids {
		var r map[string]interface{}
		if err := dao.LoadFromDB(id, &r); err != nil {
			// not loaded, get from web
			animeIDList = append(animeIDList, id)
		} else {
			resultList = append(resultList, r)
		}
	}
	// load the rest and save to db
	if len(animeIDList) != 0 {
		for _, r := range b.getPropertyByTag(animeIDList, url) {
			resultList = append(resultList, r)
			if err := dao.SaveToDB(r[`_id`].(string), r); err == nil {
				log.Println(err)
			}
		}
	}
	return resultList
}

func (b *Bangumi) getPropertyByTag(id []string, url string) []map[string]interface{} {
	collector := b.getClonedCollector()
	var response []map[string]interface{}
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	postData := make(map[string]interface{})
	postData["_ids"] = id
	postDataRaw, err := json.Marshal(postData)
	if err != nil {
		log.Fatal(err)
	}
	if err := collector.PostRaw(b.getAbsoluteURL(url), postDataRaw); err != nil {
		return nil
	}
	return response
}

func (b *Bangumi) GetLatestTorrentList() []interface{} {
	collector := b.getClonedCollector()
	response := make(map[string]interface{})
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	if err := collector.Visit(b.getAbsoluteURL(`api/torrent/latest`)); err != nil {
		return nil
	}
	return response["torrents"].([]interface{})
}

func (b *Bangumi) GetRecentTorrentList() []interface{} {
	collector := b.getClonedCollector()
	var response []interface{}
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	if err := collector.Visit(b.getAbsoluteURL(`api/bangumi/recent`)); err != nil {
		return nil
	}
	return response
}

func (b *Bangumi) DownloadTorrentFromUrl(path string, fileName string, url string) {
	downloader := b.getClonedCollector()
	downloader.OnResponse(func(r *colly.Response) {
		log.Printf("download --> %s", path+fileName)
		f, err := os.Create(path + fileName)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.Copy(f, bytes.NewReader(r.Body)); err != nil {
			log.Println(err)
		}
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	})
	if err := downloader.Visit(b.getAbsoluteURL(url)); err != nil {
		log.Fatal(err)
	}

}
