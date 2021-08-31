package module

import (
	"fmt"
	"goCrawler/dao"
	"goCrawler/qbt"
	"log"
)

type Downloader interface {
	AddTorrentToBitTorrent(path string, fileName string) error
	Contains(infoHash string) bool
	Completed(infoHash string) bool
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

func (d *DownloaderImpl) AddTorrentToBitTorrent(path string, fileName string) error {
	downloadOpts := qbt.DownloadOptions{}
	if err := d.client.DownloadFromFile(path+fileName, downloadOpts); err != nil {
		return err
	}
	return nil
}

func (d *DownloaderImpl) Contains(infoHash string) bool {
	torrentsOpts := qbt.TorrentsOptions{}
	torrentsOpts.Hashes = []string{infoHash}
	torrents, err := d.client.Torrents(torrentsOpts)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if len(torrents) > 0 {
		return true
	}
	return false
}

func (d *DownloaderImpl) Completed(infoHash string) bool {
	torrentsOpts := qbt.TorrentsOptions{}
	filter := "completed"
	torrentsOpts.Filter = &filter
	torrentsOpts.Hashes = []string{infoHash}
	torrents, err := d.client.Torrents(torrentsOpts)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if len(torrents) > 0 {
		return true
	}
	return false
}
