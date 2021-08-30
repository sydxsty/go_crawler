package dao

import (
	"log"
	"regexp"
)

type BangumiTorrentInfo struct {
	Title     string // raw title
	TorrentId string
	TeamId    string
	TagIds    []interface{}
	InfoHash  string
	Content   []interface{}
	Detail    *BangumiTorrentDetail
}

type BangumiTorrentDetail struct {
	OriginalTitle      string
	TorrentDownloadURL string
	TorrentChsName     string
	TorrentEngName     string
	TorrentJpnName     string
	TeamName           string
	Resolution         string
	Format             string
	Language           string
	Episode            string
}

func NewBangumiTorrentInfoFromMap(raw map[string]interface{}) *BangumiTorrentInfo {
	b := &BangumiTorrentInfo{}
	ok := false
	if b.Title, ok = raw["title"].(string); !ok {
		log.Fatal("cannot get title")
	}
	if b.TorrentId, ok = raw["_id"].(string); !ok {
		log.Fatal("cannot get torrent id")
	}
	if b.TeamId, ok = raw["team_id"].(string); !ok {
		log.Fatal("cannot get team tags")
	}
	if b.TagIds, ok = raw["tag_ids"].([]interface{}); !ok {
		log.Fatal("cannot get torrent tags")
	}
	if b.InfoHash, ok = raw["infoHash"].(string); !ok {
		log.Fatal("cannot get torrent info hash")
	}
	if b.Content, ok = raw["content"].([]interface{}); !ok {
		log.Fatal("cannot get torrent content")
	}
	return b
}

func (b *BangumiTorrentInfo) GenerateTorrentDetail(miscList []map[string]interface{}) {
	d := &BangumiTorrentDetail{}
	b.Detail = d
	d.OriginalTitle = b.Title
	d.TorrentDownloadURL = `/download/torrent/` + b.TorrentId + `/` + b.TorrentId + `.torrent`
	for _, misc := range miscList {
		for _, cate := range b.TagIds {
			if misc["_id"] == cate {
				switch misc["type"] {
				case "bangumi":
					translate := misc["locale"].(map[string]interface{})
					d.TorrentChsName = translate["zh_cn"].(string)
					d.TorrentEngName = translate["en"].(string)
					d.TorrentJpnName = translate["ja"].(string)
				case "team":
					d.TeamName = misc["name"].(string)
				case "resolution":
					if len(d.Resolution) != 0 {
						d.Resolution += " "
					}
					d.Resolution += misc["name"].(string)
				case "format":
					if len(d.Format) != 0 {
						d.Format += " "
					}
					d.Format += misc["name"].(string)
				case "lang":
					if len(d.Language) != 0 {
						d.Language += " "
					}
					d.Language += misc["name"].(string)
				case "misc":
					break
				default:
					log.Println(misc)
				}
			}
		}
	}
	getSeason := func() string {
		var episode string
		multi := regexp.MustCompile(`(( )|(【)|(\[))+([0-9]{1,2})(-)([0-9]{1,2})(( )|(】)|(]))`).FindAllString(d.OriginalTitle, -1)
		for _, param := range multi {
			param = regexp.MustCompile(`([0-9]{1,2})(-)([0-9]{1,2})`).FindString(param)
			episode += param
		}
		return episode
	}
	getEpisode := func() string {
		var episode string
		single := regexp.MustCompile(`(( )|\[|【|第)+([0-9]{2,3})(( )|]|】|话|話)`).FindAllString(d.OriginalTitle, -1)
		if single == nil {
			single = regexp.MustCompile(`(( )|\[|【|第)+(1)([0-9]{2,3})(( )|]|】|话|話)`).FindAllString(d.OriginalTitle, -1)
		}
		for _, param := range single {
			param = regexp.MustCompile(`([0-9]{1,4})`).FindString(param)
			episode += param
		}
		return episode
	}
	getOVA := func() string {
		var episode string
		movie := regexp.MustCompile(`(\W(剧场版)|(OVA)|(OAD)\W)`).FindAllString(d.OriginalTitle, -1)
		for _, param := range movie {
			param = regexp.MustCompile(`((剧场版)|(OVA)|(OAD))`).FindString(param)
			episode += param
		}
		return episode
	}
	d.Episode = getSeason() + getEpisode() + getOVA()
}
