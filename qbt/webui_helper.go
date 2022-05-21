package qbt

import (
	"bytes"
	"io"
	"os"
)

type WEBUIHelper interface {
	AddTorrentFromFile(path string, fileName string) error
	AddTorrentFromData(data []byte) error
	Contains(infoHash string) bool
	Completed(infoHash string) bool
}

type WEBUIHelperImpl struct {
	client *Client
}

func NewWEBUIHelper(webuiAddr string, name string, pass string) (WEBUIHelper, error) {
	w := &WEBUIHelperImpl{}
	w.client = NewClient(webuiAddr)
	loginOpts := LoginOptions{
		Username: name,
		Password: pass,
	}
	if err := w.client.Login(loginOpts); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *WEBUIHelperImpl) AddTorrentFromFile(path string, fileName string) error {
	downloadOpts := DownloadOptions{}
	return w.client.DownloadFromFile(path+fileName, downloadOpts)
}

func (w *WEBUIHelperImpl) AddTorrentFromData(data []byte) error {
	f, _ := os.Create(`tmp.torrent`)
	io.Copy(f, bytes.NewReader(data))
	f.Close()
	downloadOpts := DownloadOptions{}
	return w.client.DownloadFromFile(`tmp.torrent`, downloadOpts)
}

func (w *WEBUIHelperImpl) Contains(infoHash string) bool {
	torrentsOpts := TorrentsOptions{}
	torrentsOpts.Hashes = []string{infoHash}
	torrents, err := w.client.Torrents(torrentsOpts)
	if err != nil {
		return false
	}
	if len(torrents) > 0 {
		return true
	}
	return false
}

func (w *WEBUIHelperImpl) Completed(infoHash string) bool {
	torrentsOpts := TorrentsOptions{}
	filter := "completed"
	torrentsOpts.Filter = &filter
	torrentsOpts.Hashes = []string{infoHash}
	torrents, err := w.client.Torrents(torrentsOpts)
	if err != nil {
		return false
	}
	if len(torrents) > 0 {
		return true
	}
	return false
}
