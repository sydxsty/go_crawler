package util

import (
	"github.com/gocolly/colly/v2"
	"net/http"
	"net/url"
)

type Client interface {
	// SetRequestCallback is called before request
	SetRequestCallback(callback func(r *colly.Request))

	// SetResponseCallback is called after response
	SetResponseCallback(callback func(r *colly.Response))

	// Visit a website
	Visit(link string) error

	// Post a website
	Post(link string, requestData map[string]string) error

	// SyncVisit url after setting corresponding request and response
	SyncVisit(link string) (*colly.Response, error)

	// SyncPostRaw post raw data, can be used in posting multipart
	SyncPostRaw(link string, body []byte) (*colly.Response, error)

	SetCookies(cookies []*http.Cookie) error

	Cookies() []*http.Cookie

	// CloneBase a new ClientBase
	CloneBase() Client

	// Reset clear all req and resp func
	Reset()

	SetChild(Client)
}

type ClientBase struct {
	child     Client
	collector *colly.Collector
	domain    *url.URL
}

func NewClientBase(link string) (*ClientBase, error) {
	client := &ClientBase{}
	client.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	var err error
	client.domain, err = url.Parse(link)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SetRequestCallback is called before request
func (c *ClientBase) SetRequestCallback(callback func(r *colly.Request)) {
	c.collector.OnRequest(callback)
}

// SetResponseCallback is called after response
func (c *ClientBase) SetResponseCallback(callback func(r *colly.Response)) {
	c.collector.OnResponse(callback)
}

func (c *ClientBase) Visit(link string) error {
	defer c.child.Reset()
	return c.collector.Visit(MustGetAbsoluteURL(c.domain, link))
}

func (c *ClientBase) Post(link string, requestData map[string]string) error {
	defer c.child.Reset()
	return c.collector.Post(MustGetAbsoluteURL(c.domain, link), requestData)
}

func (c *ClientBase) SyncVisit(link string) (*colly.Response, error) {
	var resp *colly.Response
	c.child.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	if err := c.child.Visit(link); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClientBase) SyncPostRaw(link string, body []byte) (*colly.Response, error) {
	var resp *colly.Response
	c.child.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	defer c.child.Reset()
	if err := c.collector.PostRaw(MustGetAbsoluteURL(c.domain, link), body); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClientBase) SetCookies(cookies []*http.Cookie) error {
	absURL := MustGetAbsoluteURL(c.domain, "/")
	return c.collector.SetCookies(absURL, cookies)
}

func (c *ClientBase) Cookies() []*http.Cookie {
	absoluteURL := MustGetAbsoluteURL(c.domain, "/")
	return c.collector.Cookies(absoluteURL)
}

func (c *ClientBase) CloneBase() Client {
	client := &ClientBase{
		child:     nil,
		collector: c.collector,
		domain:    c.domain,
	}
	return client
}

func (c *ClientBase) Reset() {
	c.collector = c.collector.Clone()
}

func (c *ClientBase) SetChild(child Client) {
	c.child = child
}
