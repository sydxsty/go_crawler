package util

import (
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Client is used to call client
// interface isolation to ClientBase
type Client interface {
	// SetRequestCallback is called before request
	SetRequestCallback(callback func(r *colly.Request))

	// SetResponseCallback is called after response
	SetResponseCallback(callback func(r *colly.Response))

	// Visit a website
	Visit(link string) error

	// Post method, visit a website
	Post(link string, requestData map[string]string) error

	// SyncVisit url after setting corresponding request and response
	SyncVisit(link string) (*colly.Response, error)

	// SyncPostRaw post raw data, can be used in posting multipart
	SyncPostRaw(link string, body []byte) (*colly.Response, error)
}

// ClientBase is used to implement client operations
type ClientBase interface {
	Client

	SetCookies(cookies []*http.Cookie) error

	Cookies() []*http.Cookie

	// CloneBase a new ClientBase
	CloneBase(ClientBase) ClientBase

	// Reset clear all req and resp func
	// also call this func to init a client
	Reset()
}

type ClientBaseImpl struct {
	child     ClientBase
	collector *colly.Collector
	domain    *url.URL
	retry     int
	span      time.Duration // sleep between operations
}

// NewClientBase return a new ClientBase, called by constructor of child Client
func NewClientBase(child ClientBase, link string) (ClientBase, error) {
	client := &ClientBaseImpl{
		retry: 3,
		span:  3,
	}
	client.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	var err error
	client.domain, err = url.Parse(link)
	if err != nil {
		return nil, err
	}
	client.child = child
	return client, nil
}

// SetRequestCallback is called before request
func (c *ClientBaseImpl) SetRequestCallback(callback func(r *colly.Request)) {
	c.collector.OnRequest(callback)
}

// SetResponseCallback is called after response
func (c *ClientBaseImpl) SetResponseCallback(callback func(r *colly.Response)) {
	c.collector.OnResponse(callback)
}

func (c *ClientBaseImpl) callWithRetryAndDelay(callable func() error) error {
	defer c.child.Reset()
	var err error
	for i := 0; i < c.retry; i++ {
		err = callable()
		time.Sleep(time.Second * c.span)
		if err == nil {
			log.Println("Page crawled successfully.")
			break
		}
		log.Println("An error occurred when crawling pages, ", err)
	}
	if err != nil {
		return errors.Wrapf(err, "can not get link for %d times", c.retry)
	}
	return nil
}

func (c *ClientBaseImpl) Visit(link string) error {
	return c.callWithRetryAndDelay(func() error {
		return c.collector.Visit(MustGetAbsoluteURL(c.domain, link))
	})
}

func (c *ClientBaseImpl) Post(link string, requestData map[string]string) error {
	return c.callWithRetryAndDelay(func() error {
		return c.collector.Post(MustGetAbsoluteURL(c.domain, link), requestData)
	})
}

func (c *ClientBaseImpl) SyncVisit(link string) (*colly.Response, error) {
	var resp *colly.Response
	c.child.SetResponseCallback(func(r *colly.Response) {
		resp = r
	})
	if err := c.child.Visit(link); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClientBaseImpl) SyncPostRaw(link string, body []byte) (*colly.Response, error) {
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

func (c *ClientBaseImpl) SetCookies(cookies []*http.Cookie) error {
	absURL := MustGetAbsoluteURL(c.domain, "/")
	return c.collector.SetCookies(absURL, cookies)
}

func (c *ClientBaseImpl) Cookies() []*http.Cookie {
	absoluteURL := MustGetAbsoluteURL(c.domain, "/")
	return c.collector.Cookies(absoluteURL)
}

func (c *ClientBaseImpl) CloneBase(child ClientBase) ClientBase {
	client := &ClientBaseImpl{
		child:     child,
		collector: c.collector,
		domain:    c.domain,
		retry:     c.retry,
		span:      c.span,
	}
	return client
}

func (c *ClientBaseImpl) Reset() {
	c.collector = c.collector.Clone()
}
