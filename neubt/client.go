package neubt

import (
	"crawler/storage"
	"crawler/util"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"net/http"
	"net/url"
)

type Client interface {
	Login(username string, password string) error
	LoadCookie(cookiePath string) error
	SaveCookie(cookiePath string) error
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

type ClientImpl struct {
	collector *colly.Collector
	db        storage.KVStorage
	domain    *url.URL
}

func NewClient(db storage.KVStorage) (Client, error) {
	client := &ClientImpl{}
	client.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	client.domain, _ = url.Parse("http://[2001:da8:9000::232]")
	client.db = db
	client.Reset()
	return client, nil
}

func (c *ClientImpl) Login(username string, password string) error {
	absURL := util.MustGetAbsoluteURL(c.domain, "member.php?mod=logging&action=login&loginsubmit=yes&inajax=1")
	defer c.Reset()
	return c.collector.Post(absURL, map[string]string{
		"username":   username,
		"password":   password,
		"questionid": "0",
		"answer":     "",
	})
}

func (c *ClientImpl) LoadCookie(cookiePath string) error {
	var cookie []*http.Cookie
	raw, err := c.db.GetRaw(cookiePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &cookie)
	if err != nil {
		return err
	}
	absURL := util.MustGetAbsoluteURL(c.domain, "/")
	return c.collector.SetCookies(absURL, cookie)
}

func (c *ClientImpl) SaveCookie(cookiePath string) error {
	absoluteURL := util.MustGetAbsoluteURL(c.domain, "/")
	raw, err := json.Marshal(c.collector.Cookies(absoluteURL))
	if err != nil {
		return err
	}
	return c.db.PutRaw(cookiePath, raw)
}

// SetRequestCallback is called before request
func (c *ClientImpl) SetRequestCallback(callback func(r *colly.Request)) {
	c.collector.OnRequest(callback)
}

// SetResponseCallback is called after response
func (c *ClientImpl) SetResponseCallback(callback func(r *colly.Response)) {
	c.collector.OnResponse(callback)
}

func (c *ClientImpl) visit(link string) error {
	defer c.Reset()
	return c.collector.Visit(util.MustGetAbsoluteURL(c.domain, link))
}

func (c *ClientImpl) SyncVisit(link string) (*colly.Response, error) {
	var resp *colly.Response
	c.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	if err := c.visit(link); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClientImpl) SyncPostRaw(link string, body []byte) (*colly.Response, error) {
	var resp *colly.Response
	c.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	defer c.Reset()
	if err := c.collector.PostRaw(util.MustGetAbsoluteURL(c.domain, link), body); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClientImpl) Clone() Client {
	client := &ClientImpl{
		collector: c.collector,
		db:        c.db,
		domain:    c.domain,
	}
	client.Reset()
	return client
}

func (c *ClientImpl) Reset() {
	c.collector = c.collector.Clone()
	c.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Host", c.domain.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Origin", util.MustGetAbsoluteURL(c.domain, "/"))
		r.Headers.Set("Referer", util.MustGetAbsoluteURL(c.domain, "plugin.php?id=neubt_resourceindex"))
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
}
