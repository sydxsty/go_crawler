package crawler

import (
	"crawler/bangumi"
	"log"
)

type Downloader interface {
	DownloadTorrentFromUrl(link string) ([]byte, error)
}

type DownloaderImpl struct {
	client bangumi.Client
}

func NewDownloader(client bangumi.Client) Downloader {
	return &DownloaderImpl{
		client: client.Clone(),
	}
}

func (d *DownloaderImpl) DownloadTorrentFromUrl(link string) ([]byte, error) {
	log.Println("download torrent from ", link)
	resp, err := d.client.SyncVisit(link)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
