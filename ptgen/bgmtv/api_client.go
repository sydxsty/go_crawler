package bgmtv

import (
	"crawler/util"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"log"
	"time"
)

type APIClientImpl struct {
	util.ClientBase // parent
	retry           int
	span            time.Duration
}

func NewAPIClient() (Client, error) {
	client := &APIClientImpl{
		retry: 3,
		span:  3,
	}
	var err error
	client.ClientBase, err = util.NewClientBase(client, "https://api.bgm.tv")
	if err != nil {
		return nil, err
	}
	client.Reset()
	return client, nil
}

func (c *APIClientImpl) SyncVisit(link string) (*colly.Response, error) {
	r, err := c.ClientBase.SyncVisit(link)
	for i := 0; i < c.retry-1 && err != nil; i++ {
		r, err = c.ClientBase.SyncVisit(link)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "can not get link for %d times", c.retry)
	}
	time.Sleep(time.Second * c.span)
	log.Println(c.Cookies())
	return r, nil
}

func (c *APIClientImpl) Clone() Client {
	client := &APIClientImpl{
		retry: c.retry,
		span:  c.span,
	}
	client.ClientBase = c.ClientBase.CloneBase(client)
	client.Reset()
	return client
}

func (c *APIClientImpl) Reset() {
	c.ClientBase.Reset()
	c.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Host", "api.bgm.tv")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36`)
		r.Headers.Set("Accept", "application/json")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Accept-Encoding", "deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
}
