package bangumi

import (
	"crawler/util"
	"github.com/gocolly/colly/v2"
	"net/url"
)

type Client interface {
	// SetRequestCallback is called before request
	SetRequestCallback(callback func(r *colly.Request))

	// SetResponseCallback is called after response
	SetResponseCallback(callback func(r *colly.Response))

	// SyncVisit url after setting corresponding request and response
	SyncVisit(link string) (*colly.Response, error)

	// SyncPostRaw post raw data, can be used in posting multipart
	SyncPostRaw(link string, body []byte) (*colly.Response, error)

	// Clone a new Client
	Clone() Client

	// Reset clear all req and resp func
	Reset()
}

type BangumiClient struct {
	collector *colly.Collector
	domain    *url.URL
}

func NewClient() (Client, error) {
	client := &BangumiClient{}
	client.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	client.domain, _ = url.Parse("https://bangumi.moe")
	client.Reset()
	return client, nil
}

// SetRequestCallback is called before request
func (b *BangumiClient) SetRequestCallback(callback func(r *colly.Request)) {
	b.collector.OnRequest(callback)
}

// SetResponseCallback is called after response
func (b *BangumiClient) SetResponseCallback(callback func(r *colly.Response)) {
	b.collector.OnResponse(callback)
}

func (b *BangumiClient) visit(link string) error {
	defer b.Reset()
	return b.collector.Visit(util.MustGetAbsoluteURL(b.domain, link))
}

func (b *BangumiClient) SyncVisit(link string) (*colly.Response, error) {
	var resp *colly.Response
	b.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	if err := b.visit(link); err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BangumiClient) SyncPostRaw(link string, body []byte) (*colly.Response, error) {
	var resp *colly.Response
	b.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	defer b.Reset()
	if err := b.collector.PostRaw(util.MustGetAbsoluteURL(b.domain, link), body); err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BangumiClient) Clone() Client {
	client := &BangumiClient{
		collector: b.collector,
		domain:    b.domain,
	}
	client.Reset()
	return client
}

func (b *BangumiClient) Reset() {
	b.collector = b.collector.Clone()
	b.SetRequestCallback(func(r *colly.Request) {
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
}
