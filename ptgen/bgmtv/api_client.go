package bgmtv

import (
	"crawler/util"
	"github.com/gocolly/colly/v2"
)

type APIClientImpl struct {
	util.ClientBase // parent
}

func NewAPIClient() (Client, error) {
	client := &APIClientImpl{}
	var err error
	client.ClientBase, err = util.NewClientBase(client, "https://api.bgm.tv")
	if err != nil {
		return nil, err
	}
	client.Reset()
	return client, nil
}

func (c *APIClientImpl) Clone() Client {
	client := &APIClientImpl{}
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
