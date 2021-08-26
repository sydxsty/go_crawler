package module

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"goCrawler/qbt"
	"io"
	"log"
	"os"
)

type Downloader interface {
	DownloadTorrent(link string, fileName string, path string) error
	AddTorrentToBitTorrent(fileName string, path string) error
}

type DownloaderImpl struct {
	scraperModuleImpl
	fileName string
	client   *qbt.Client
}

func NewDownloader() Downloader {
	d := &DownloaderImpl{}
	d.init()
	d.client = qbt.NewClient(dao.YAMLConfig.QBAddr)
	loginOpts := qbt.LoginOptions{
		Username: dao.YAMLConfig.QBUsername,
		Password: dao.YAMLConfig.QBPassword,
	}
	if err := d.client.Login(loginOpts); err != nil {
		log.Println(err)
	}
	return d
}

func (d *DownloaderImpl) AddTorrentToBitTorrent(fileName string, path string) error {
	downloadOpts := qbt.DownloadOptions{}
	if err := d.client.DownloadFromFile(path+fileName, downloadOpts); err != nil {
		return err
	}
	// ******************
	// GET ALL TORRENTS *
	// ******************
	torrentsOpts := qbt.TorrentsOptions{}
	filter := "inactive"
	sort := "name"
	hash := "d739f78a12b241ba62719b1064701ffbb45498a8"
	torrentsOpts.Filter = &filter
	torrentsOpts.Sort = &sort
	torrentsOpts.Hashes = []string{hash}
	torrents, err := d.client.Torrents(torrentsOpts)
	if err != nil {
		fmt.Println("[-] Get torrent list")
		fmt.Println(err)
	} else {
		fmt.Println("[+] Get torrent list")
		if len(torrents) > 0 {
			spew.Dump(torrents[0])
		} else {
			fmt.Println("No torrents found")
		}
	}
	return nil
}

func (d *DownloaderImpl) DownloadTorrent(link string, fileName string, path string) error {
	collector := d.getClonedCollector()
	var actualUrl string
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		actualUrl = node.GetInnerNode(`//*[@class="alert_btnleft"]/a/@href`).GetString()
	})
	if err := collector.Visit(d.getAbsoluteURL(link)); err != nil {
		return err
	}
	if actualUrl != "" {
		downloader := d.getClonedCollector()
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
		if err := downloader.Visit(d.getAbsoluteURL(actualUrl)); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
