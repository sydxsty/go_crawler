package dao

import (
	"crawler/neubt/html"
	"crawler/storage"
	"errors"
	"log"
	"strconv"
	"strings"
)

func NodeListToTorrentInfoList(db storage.KVStorage, nodeList []*html.NodeDecorator) ([]*TorrentInfo, error) {
	var tl []*TorrentInfo
	for _, node := range nodeList {
		rawInfo := nodeToTorrentInfo(node)
		info, err := newAggregatedTorrentInfo(db, rawInfo)
		if err != nil {
			return nil, err
		}
		tl = append(tl, info)
	}
	return tl, nil
}

func nodeToTorrentInfo(node *html.NodeDecorator) *rawTorrentInfo {
	raw := &rawTorrentInfo{}
	// generate detailed information of a torrent
	raw.Title, _ = node.GetInnerString(`.//a[@class="s xst"]`)
	raw.Link, _ = node.GetInnerString(`.//a[@class="s xst"]/@href`)
	raw.Size, _ = node.GetInnerString(`./td[3]`)
	raw.Discount, _ = node.GetInnerString(`.//a[@class="s xst"]/../img[last()]/@src`)
	raw.Signal, _ = node.GetInnerString(`./td[2]/img/@src`)
	raw.Uploader, _ = node.GetInnerString(`./td[last()]//a`)
	raw.UploaderPage, _ = node.GetInnerString(`./td[last()]//a/@href`)
	return raw
}

func newAggregatedTorrentInfo(db storage.KVStorage, r *rawTorrentInfo) (*TorrentInfo, error) {
	t := &TorrentInfo{}
	if err := db.Get(r.Link, t); err != nil {
		// key not exist
		if err := t.loadFromRawTorrentInfo(r); err != nil {
			// phase error
			return nil, err // error
		}
		return t, nil
	}
	// get key successful, merge them
	if err := t.updateFromRawTorrentInfo(r); err != nil {
		return nil, err // db only
	}
	if err := db.Put(t.Link, t); err != nil {
		return nil, err // db only
	}
	return t, nil // db + raw
}

func (t *TorrentInfo) loadFromRawTorrentInfo(r *rawTorrentInfo) error {
	return t.applyFilter(
		titleFilter(r),
		linkFilter(r),
		sizeFilter(r),
		discountFilter(r),
		hasCrawledFilter(r),
	)
}

func (t *TorrentInfo) updateFromRawTorrentInfo(r *rawTorrentInfo) error {
	return t.applyFilter(
		titleFilter(r),
		discountFilter(r),
	)
}

type filterOption func(*TorrentInfo) error

func (t *TorrentInfo) applyFilter(options ...filterOption) error {
	for _, f := range options {
		err := f(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func titleFilter(r *rawTorrentInfo) filterOption {
	return func(p *TorrentInfo) error {
		if len(r.Title) == 0 {
			return errors.New("torrent title is empty")
		}
		p.Title = r.Title
		return nil
	}
}

func linkFilter(r *rawTorrentInfo) filterOption {
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

func sizeFilter(r *rawTorrentInfo) filterOption {
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

func discountFilter(r *rawTorrentInfo) filterOption {
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

func hasCrawledFilter(r *rawTorrentInfo) filterOption {
	return func(p *TorrentInfo) error {
		return nil
	}
}
