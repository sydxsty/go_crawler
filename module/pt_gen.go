package module

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"log"
	"net/url"
)

type PTGen struct {
	scraperModuleImpl
}

func NewPTGen() *PTGen {
	p := &PTGen{}
	p.init()
	p.domain, _ = url.Parse("https://api.douban.workers.dev")
	return p
}

func (p *PTGen) getClonedCollector() *colly.Collector {
	clone := p.collector.Clone()
	clone.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "api.douban.workers.dev")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`)
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Origin", `https://api.douban.workers.dev`)
		r.Headers.Set("Referer", `https://api.douban.workers.dev/`)
		r.Headers.Set("Accept-Encoding", "deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
	return clone
}

func (p *PTGen) GetBangumiLinkByName(name string) map[string]string {
	if name == "" {
		return nil
	}
	collector := p.getClonedCollector()
	var response map[string]interface{}
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	if err := collector.Visit(p.getAbsoluteURL(`/?search=` + name + `&source=bangumi`)); err != nil {
		return nil
	}
	respData, ok := response["data"].([]interface{})
	if !ok {
		return nil
	}
	linkMap := make(map[string]string)
	for _, node := range respData {
		unmarshalNode := node.(map[string]interface{})
		if unmarshalNode["subtype"].(string) == "动画/二次元番" {
			linkMap[unmarshalNode["title"].(string)] = unmarshalNode["link"].(string)
			break
		}
	}
	return linkMap
}

func (p *PTGen) GetBangumiDetailByLink(link string) map[string]interface{} {
	var response map[string]interface{}
	// load from buffer
	if err := dao.LoadFromDB("pt_gen_link_"+link, &response); err == nil {
		return response
	}
	collector := p.getClonedCollector()
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &response); err != nil {
			log.Fatal(err)
		}
	})
	if err := collector.Visit(p.getAbsoluteURL(`/?url=` + link)); err != nil {
		return nil
	}
	if response != nil {
		// save to buffer
		if err := dao.SaveToDB("pt_gen_link_"+link, response); err != nil {
			return nil
		}
	}
	return response
}
