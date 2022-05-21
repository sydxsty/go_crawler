package crawler

import (
	"crawler/bangumi"
	"log"
)

type Downloader interface {
	DownloadTorrentFromUrl(url string) ([]byte, error)
}

type DownloaderImpl struct {
	client bangumi.Client
}

func NewDownloader(client bangumi.Client) Downloader {
	return &DownloaderImpl{
		client: client.Clone(),
	}
}

func (d *DownloaderImpl) DownloadTorrentFromUrl(url string) ([]byte, error) {
	log.Println("download torrent from ", url)
	resp, err := d.client.SyncVisit(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
