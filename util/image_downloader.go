package util

import (
	"github.com/gocolly/colly/v2"
	"net/url"
	"regexp"
)

type ImageDownloader interface {
	// Download return file data, file type
	Download() ([]byte, string, error)
}

type ImageDownloaderImpl struct {
	ClientBase // parent
	domain     *url.URL
	link       string
}

func NewImageDownloader(link string) (ImageDownloader, error) {
	client := &ImageDownloaderImpl{}
	client.link = link
	var err error
	client.domain, err = url.Parse(link)
	if err != nil {
		return nil, err
	}
	client.ClientBase, err = NewClientBase(client, link)
	if err != nil {
		return nil, err
	}
	client.Reset()
	return client, nil
}

func (c *ImageDownloaderImpl) Reset() {
	c.ClientBase.Reset()
	c.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Host", c.domain.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`)
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "text/plain;charset=UTF-8")
		r.Headers.Set("Origin", MustGetAbsoluteURL(c.domain, ""))
		r.Headers.Set("Referer", MustGetAbsoluteURL(c.domain, "/"))
		r.Headers.Set("Accept-Encoding", "deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
}

func (c *ImageDownloaderImpl) Download() ([]byte, string, error) {
	resp, err := c.SyncVisit(c.link)
	if err != nil {
		return nil, "", err
	}
	fileType := regexp.MustCompile(`[^.]\w*$`).FindString(c.link)
	return resp.Body, fileType, err
}
