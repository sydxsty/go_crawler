package module

import (
	"bytes"
	"errors"
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

type ScraperModule interface {
	Login(username string, password string)
	SaveCookie() error
	LoadCookie() error
	DownloadTorrentFromNestedURL(link string, fileName string, path string) error
	DownloadTorrentFromResponseBody(r *colly.Response, fileName string, path string) error
}

type scraperModuleImpl struct {
	domain    *url.URL
	collector *colly.Collector
}

var basicModule *scraperModuleImpl

func init() {
	basicModule = &scraperModuleImpl{}
	d, _ := url.Parse("http://[2001:da8:9000::232]")
	basicModule.domain = d
	basicModule.collector = colly.NewCollector()
	if err := basicModule.collector.SetProxy("http://127.0.0.1:8888"); err != nil {
		log.Fatal(err)
	}
	basicModule.collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"
}

func (b *scraperModuleImpl) init() {
	b.collector = basicModule.collector
	b.domain = basicModule.domain
	b.collector.AllowURLRevisit = true
}

func (b *scraperModuleImpl) Login(username string, password string) {
	if err := b.collector.Post(b.getAbsoluteURL("member.php?mod=logging&action=login&loginsubmit=yes&inajax=1"),
		map[string]string{"username": username, "password": password, "questionid": "0", "answer": ""}); err != nil {
		log.Fatal(err)
	}
}

func (b *scraperModuleImpl) SaveCookie() error {
	return dao.SaveCookieToDB(b.collector.Cookies(b.getAbsoluteURL("/")))
}

func (b *scraperModuleImpl) LoadCookie() error {
	cookie, err := dao.LoadCookieFromDB()
	if err != nil {
		return err
	}
	if err := b.collector.SetCookies(b.getAbsoluteURL("/"), cookie); err != nil {
		return err
	}
	return nil
}

func (b *scraperModuleImpl) getAbsoluteURL(u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	absURL, err := b.domain.Parse(u)
	if err != nil {
		return ""
	}
	absURL.Fragment = ""
	return absURL.String()
}

func (b *scraperModuleImpl) getClonedCollector() *colly.Collector {
	clone := b.collector.Clone()
	clone.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", b.domain.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Origin", b.getAbsoluteURL("/"))
		r.Headers.Set("Referer", b.getAbsoluteURL("plugin.php?id=neubt_resourceindex"))
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
	return clone
}

func (b *scraperModuleImpl) DownloadTorrentFromNestedURL(link string, fileName string, path string) error {
	collector := b.getClonedCollector()
	var response *colly.Response
	collector.OnResponse(func(r *colly.Response) {
		response = r
	})
	if err := collector.Visit(b.getAbsoluteURL(link)); err != nil {
		return err
	}
	return b.DownloadTorrentFromResponseBody(response, fileName, path)
}

func (b *scraperModuleImpl) DownloadTorrentFromResponseBody(r *colly.Response, fileName string, path string) error {
	var actualUrl string
	node, err := NewNodeFromBytes(r.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	actualUrl = node.GetInnerNode(`//*[@class="alert_btnleft"]/a/@href`).GetString()

	if actualUrl != "" {
		downloader := b.getClonedCollector()
		downloader.OnResponse(func(r *colly.Response) {
			log.Printf("download --> %s", path+fileName)
			f, err := os.Create(path + fileName)
			if err != nil {
				log.Println(err)
				return
			}
			if _, err := io.Copy(f, bytes.NewReader(r.Body)); err != nil {
				log.Println(err)
			}
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		})
		if err := downloader.Visit(b.getAbsoluteURL(actualUrl)); err != nil {
			log.Fatal(err)
		}
		return nil
	} else {
		return errors.New("url not found")
	}
}
