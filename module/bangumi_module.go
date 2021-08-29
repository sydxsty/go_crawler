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

func (b *Bangumi) GetAnimeNameByTag(id ...string) map[interface{}]interface{} {
	collector := b.getClonedCollector()
	var response map[interface{}]interface{}
	collector.OnResponse(func(r *colly.Response) {
		// json data
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	rawID, _ := json.Marshal(id)
	if err := collector.Post(b.getAbsoluteURL(`api/tag/fetch`), map[string]string{"_ids": string(rawID)}); err != nil {
		return nil
	}
	return response
}

