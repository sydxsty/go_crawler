package ptgen

import (
	"crawler/util"
	"github.com/gocolly/colly/v2"
)

type Client interface {
	util.Client
	Clone() Client
	Reset()
}

type ClientImpl struct {
	util.ClientBase // parent
}

func NewClient() (Client, error) {
	client := &ClientImpl{}
	var err error
	client.ClientBase, err = util.NewClientBase("https://ptgen.frfx.workers.dev")
	if err != nil {
		return nil, err
	}
	client.SetChild(client)
	client.Reset()
	return client, nil
}

func (c *ClientImpl) Clone() Client {
	client := &ClientImpl{
		ClientBase: c.ClientBase.CloneBase(),
	}
	client.SetChild(client)
	client.Reset()
	return client
}

func (c *ClientImpl) Reset() {
	c.ClientBase.Reset()
	c.SetRequestCallback(func(r *colly.Request) {
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
