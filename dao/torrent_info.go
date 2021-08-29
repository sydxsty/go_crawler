package dao

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RawTorrentInfo struct {
	Title        string
	Link         string
	Size         string
	Discount     string
	Signal       string
	Uploader     string
	UploaderPage string
}

type TorrentInfo struct {
	Title     string
	Link      string
	TorrentID int
	Size      float64 // MB
	Discount  int     // 0 to 100
	Crawled   bool
}

type FilterOption func(*TorrentInfo) error

func NewAggregatedTorrentInfo(r *RawTorrentInfo) *TorrentInfo {
	t := &TorrentInfo{}
	if err := t.LoadFromDB(r.Link); err != nil {
		if err := t.LoadFromRawTorrentInfo(r); err != nil {
			return nil // error
		}
		if err := t.SaveToDB(); err != nil {
			log.Println(err)
		}
		return t // raw only
	}
	if err := t.UpdateFromRawTorrentInfo(r); err != nil {
		return t // db only
	}
	if err := t.SaveToDB(); err != nil {
		log.Println(err)
	}
	return t // db + raw
}

func (t *TorrentInfo) LoadFromRawTorrentInfo(r *RawTorrentInfo) error {
	return t.applyFilter(
		titleFilter(r),
		linkFilter(r),
		sizeFilter(r),
		discountFilter(r),
		hasCrawledFilter(r),
	)
}

func (t *TorrentInfo) UpdateFromRawTorrentInfo(r *RawTorrentInfo) error {
	return t.applyFilter(
		titleFilter(r),
		discountFilter(r),
	)
}

func (t *TorrentInfo) applyFilter(options ...FilterOption) error {
	for _, f := range options {
		err := f(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TorrentInfo) LoadFromDB(link string) error {
	value, err := TorrentInfoDBHandle.Get([]byte(link), nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(value, t); err != nil {
		return err
	}
	return nil
}

func (t *TorrentInfo) SaveToDB() error {
	raw, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := TorrentInfoDBHandle.Put([]byte(t.Link), raw, nil); err != nil {
		return err
	}
	return nil
}

func LoadCookieFromDB() ([]*http.Cookie, error) {
	var cookie []*http.Cookie
	raw, err := TorrentInfoDBHandle.Get([]byte(YAMLConfig.CookiePath), nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &cookie); err != nil {
		return nil, err
	}
	return cookie, nil
}

func SaveCookieToDB(cookie []*http.Cookie) error {
	raw, err := json.Marshal(cookie)
	if err != nil {
		return err
	}
	if err := TorrentInfoDBHandle.Put([]byte(YAMLConfig.CookiePath), raw, nil); err != nil {
		return err
	}
	return nil
}

func titleFilter(r *RawTorrentInfo) FilterOption {
	return func(p *TorrentInfo) error {
		if len(r.Title) == 0 {
			return errors.New("torrent title is empty")
		}
		p.Title = r.Title
		return nil
	}
}

func linkFilter(r *RawTorrentInfo) FilterOption {
	return func(p *TorrentInfo) error {
		p.Link = r.Link
		result := strings.Split(r.Link, "-")
		if len(result) == 4 {
			p.TorrentID, _ = strconv.Atoi(result[1])
			return nil
		}
		// we must get link, so return error
		return errors.New("can not get torrent link")
	}
}

func sizeFilter(r *RawTorrentInfo) FilterOption {
	return func(p *TorrentInfo) error {
		result := strings.Split(r.Size, " ")
		if len(result) != 2 {
			log.Println("can not get torrent size")
			return nil
		}
		p.Size, _ = strconv.ParseFloat(result[0], 64)
		switch result[1] {
		case "GB":
			p.Size *= 1024
		case "TB":
			p.Size *= 1024 * 1024
		case "MB":
			break
		default:
			log.Println("can not get torrent size")
			return nil
		}
		return nil
	}
}

func discountFilter(r *RawTorrentInfo) FilterOption {
	return func(p *TorrentInfo) error {
		if strings.Contains(r.Discount, "free") {
			p.Discount = 100
		} else if strings.Contains(r.Discount, "dl50") {
			p.Discount = 50
		} else {
			p.Discount = 0
		}
		return nil
	}
}

func hasCrawledFilter(r *RawTorrentInfo) FilterOption {
	return func(p *TorrentInfo) error {
		return nil
	}
}
