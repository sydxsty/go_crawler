package crawler

import (
	"crawler/neubt"
	"crawler/neubt/html"
	"errors"
)

// Downloader download torrent from url like http://bt.neu6.edu.cn/forum.php?mod=attachment&aid=XXX
type Downloader interface {
	// DownloadFromNestedURL return the name and the bytecode of selected torrent
	DownloadFromNestedURL(link string) ([]byte, error)
}

type DownloaderImpl struct {
	client neubt.Client
}

func NewDownloader(client neubt.Client) Downloader {
	d := &DownloaderImpl{
		client: client.Clone(),
	}
	return d
}

func (d *DownloaderImpl) DownloadFromNestedURL(link string) ([]byte, error) {
	resp, err := d.client.SyncVisit(link)
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}
	actualUrl, err := node.GetInnerString(`//*[@class="alert_btnleft"]/a/@href`)
	if err != nil {
		return nil, err
	}
	if actualUrl == "" {
		return nil, errors.New("torrent url is empty")
	}
	return d.downloadFromDirectURL(actualUrl)
}

func (d *DownloaderImpl) downloadFromDirectURL(link string) ([]byte, error) {
	resp, err := d.client.SyncVisit(link)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
