package dao

import (
	bgm "crawler/bangumi/anime_control"
	"errors"
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

func NewTorrentInfoFromMap(raw map[string]interface{}) (*BangumiTorrentInfo, error) {
	b := &BangumiTorrentInfo{}
	ok := false
	if b.Title, ok = raw["title"].(string); !ok {
		return nil, errors.New("cannot get title")
	}
	if b.TorrentId, ok = raw["_id"].(string); !ok {
		return nil, errors.New("cannot get torrent id")
	}
	if b.TeamId, ok = raw["team_id"].(string); !ok {
		b.TeamId = "Unknown"
	}
	if b.TagIds, ok = raw["tag_ids"].([]interface{}); !ok {
		return nil, errors.New("cannot get torrent tags")
	}
	if b.InfoHash, ok = raw["infoHash"].(string); !ok {
		return nil, errors.New("cannot get torrent info hash")
	}
	if b.Content, ok = raw["content"].([]interface{}); !ok {
		return nil, errors.New("cannot get torrent content")
	}
	return b, nil
}

func (b *BangumiTorrentInfo) InitTorrentDetail(miscList []map[string]interface{}) {
	bgmFilter := bgm.NewBangumiFilter()
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
					d.TorrentChsName, _ = translate["zh_cn"].(string)
					d.TorrentEngName, _ = translate["en"].(string)
					d.TorrentJpnName, _ = translate["ja"].(string)
				case "lang":
					if len(d.Language) != 0 {
						d.Language += " "
					}
					d.Language += misc["name"].(string)
				default:
					break
				}
			}
		}
	}
	// TODO: get language detail
	// TODO: get bangumi anime name detail
	d.Resolution = getString(bgmFilter.GetResolution(b.Title))
	d.Format = getString(bgmFilter.GetMediaInfo(b.Title))
	d.TeamName = getString(bgmFilter.GetTeam(b.Title))
	d.Episode = getString([]string{getString(bgmFilter.GetMovieType(b.Title)), bgmFilter.GetSingleEpisode(b.Title), bgmFilter.GetMultiEpisode(b.Title)})
}

func getString(strList []string) string {
	res := ""
	for _, str := range strList {
		if str == "" {
			continue
		}
		if len(res) != 0 {
			res += " "
		}
		res += str
	}
	return res
}

func (b *BangumiTorrentInfo) GetCHNName() (string, error) {
	if b.Detail == nil {
		return "", errors.New("torrent detail is not init")
	}
	if len(b.Detail.TorrentChsName) != 0 {
		return b.Detail.TorrentChsName, nil
	}
	if len(b.Detail.TorrentEngName) != 0 {
		return b.Detail.TorrentEngName, nil
	}
	return "", errors.New("no valid torrent chs name found")
}
