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
	util.Client
	Login(username string, password string) error
	LoadCookie(cookiePath string) error
	SaveCookie(cookiePath string) error
	Clone() Client
	Reset()
}

type ClientImpl struct {
	util.Client // parent
	db          storage.KVStorage
}

func NewClient(db storage.KVStorage) (Client, error) {
	client := &ClientImpl{}
	var err error
	client.Client, err = util.NewClientBase("http://[2001:da8:9000::232]")
	if err != nil {
		return nil, err
	}
	client.db = db
	client.SetChild(client)
	client.Reset()
	return client, nil
}

func (c *ClientImpl) Login(username string, password string) error {
	return c.Post("member.php?mod=logging&action=login&loginsubmit=yes&inajax=1", map[string]string{
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
	return c.SetCookies(cookie)
}

func (c *ClientImpl) SaveCookie(cookiePath string) error {
	raw, err := json.Marshal(c.Cookies())
	if err != nil {
		return err
	}
	return c.db.PutRaw(cookiePath, raw)
}

func (c *ClientImpl) Clone() Client {
	client := &ClientImpl{
		Client: c.Client.CloneBase(),
		db:     c.db,
	}
	client.SetChild(client)
	client.Reset()
	return client
}

func (c *ClientImpl) Reset() {
	c.Client.Reset()
	domain, _ := url.Parse(`http://bt.neu6.edu.cn`)
	c.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Host", domain.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Origin", util.MustGetAbsoluteURL(domain, "/"))
		r.Headers.Set("Referer", util.MustGetAbsoluteURL(domain, "plugin.php?id=neubt_resourceindex"))
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
}
