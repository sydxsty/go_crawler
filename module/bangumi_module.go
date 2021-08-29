package module

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
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

func (b *Bangumi) GetAnimeNameByTag(id []string) []interface{} {
	collector := colly.NewCollector()
	var response []interface{}
	collector.OnRequest(func(r *colly.Request) {
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
	collector.OnResponse(func(r *colly.Response) {
		log.Println(r.Request.Headers)
		// json data
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
	if err := collector.PostRaw(b.getAbsoluteURL(`api/tag/fetch`), postDataRaw); err != nil {
		return nil
	}
	return response
}
