package ptgen

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

type PTGenClient struct {
	collector *colly.Collector
	domain    *url.URL
}

func NewClient() (Client, error) {
	client := &PTGenClient{}
	client.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	client.domain, _ = url.Parse("https://ptgen.frfx.workers.dev")
	client.Reset()
	return client, nil
}

// SetRequestCallback is called before request
func (p *PTGenClient) SetRequestCallback(callback func(r *colly.Request)) {
	p.collector.OnRequest(callback)
}

// SetResponseCallback is called after response
func (p *PTGenClient) SetResponseCallback(callback func(r *colly.Response)) {
	p.collector.OnResponse(callback)
}

func (p *PTGenClient) visit(link string) error {
	defer p.Reset()
	return p.collector.Visit(util.MustGetAbsoluteURL(p.domain, link))
}

func (p *PTGenClient) SyncVisit(link string) (*colly.Response, error) {
	var resp *colly.Response
	p.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	if err := p.visit(link); err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PTGenClient) SyncPostRaw(link string, body []byte) (*colly.Response, error) {
	var resp *colly.Response
	p.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	defer p.Reset()
	if err := p.collector.PostRaw(util.MustGetAbsoluteURL(p.domain, link), body); err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PTGenClient) Clone() Client {
	client := &PTGenClient{
		collector: p.collector,
		domain:    p.domain,
	}
	client.Reset()
	return client
}

func (p *PTGenClient) Reset() {
	p.collector = p.collector.Clone()
	p.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Host", "ptgen.frfx.workers.dev")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`)
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Origin", `https://ptgen.frfx.workers.dev`)
		r.Headers.Set("Referer", `https://ptgen.frfx.workers.dev/`)
		r.Headers.Set("Accept-Encoding", "deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
}
