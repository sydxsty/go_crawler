package dao

import (
	bgm "crawler/bangumi/anime_control"
	"crawler/util"
	"errors"
	"log"
)

type BangumiTorrentInfo struct {
	bgmFilter *bgm.BangumiFilter
	Title     string // raw title
	TorrentId string
	TeamId    string
	TagIds    []interface{}
	InfoHash  string
	content   []interface{}
	detail    *BangumiTorrentDetail
}

type BangumiTorrentDetail struct {
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
	b.bgmFilter = bgm.NewBangumiFilter()
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
	if b.content, ok = raw["content"].([]interface{}); !ok {
		return nil, errors.New("cannot get torrent content")
	}
	return b, nil
}

func (b *BangumiTorrentInfo) InitTorrentDetail(miscList []map[string]interface{}) {
	d := &BangumiTorrentDetail{}
	b.detail = d
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
	d.Resolution = getString(b.bgmFilter.GetResolution(b.Title))
	d.Format = getString(b.bgmFilter.GetMediaInfo(b.Title))
	d.TeamName = getString(b.bgmFilter.GetTeam(b.Title))
	d.Episode = getString([]string{
		getString(b.bgmFilter.GetMovieType(b.Title)),
		b.bgmFilter.GetSingleEpisode(b.Title),
		b.bgmFilter.GetMultiEpisode(b.Title),
	})
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

func (b *BangumiTorrentInfo) GetTorrentDownloadURL() (string, error) {
	if b.detail == nil {
		return "", errors.New("not init")
	}
	if len(b.detail.TorrentChsName) != 0 {
		return b.detail.TorrentDownloadURL, nil
	}
	return "", errors.New("can not get torrent url")
}

func (b *BangumiTorrentInfo) MustGetCHSName() string {
	if b.detail == nil {
		return ""
	}
	if len(b.detail.TorrentChsName) != 0 {
		return b.detail.TorrentChsName
	}
	if len(b.detail.TorrentEngName) != 0 {
		return b.detail.TorrentEngName
	}
	return ""
}

func (b *BangumiTorrentInfo) MustGetENGName() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.TorrentEngName
}

func (b *BangumiTorrentInfo) MustGetJPNName() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.TorrentJpnName
}

func (b *BangumiTorrentInfo) MustGetTeamName() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.TeamName
}

func (b *BangumiTorrentInfo) MustGetResolution() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.Resolution
}

func (b *BangumiTorrentInfo) MustGetEpisode() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.Episode
}

func (b *BangumiTorrentInfo) MustGetFormat() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.Format
}

func (b *BangumiTorrentInfo) MustGetLanguage() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.Language
}

// SetReleaseCHSName change the chs name to release version, the chs name may not be searchable in ptgen
func (b *BangumiTorrentInfo) SetReleaseCHSName(name string) {
	if len(name) != 0 {
		log.Println("set torrent chsName to: ", name)
		b.detail.TorrentChsName = name
	}
	target := b.bgmFilter.GetSeasonType(b.detail.TorrentChsName)
	if len(target) == 0 { // the season is empty
		b.detail.TorrentChsName += getString(b.bgmFilter.GetSeasonType(b.Title))
	}
}

func (b *BangumiTorrentInfo) SetCHSName(name string) {
	b.detail.TorrentChsName = name
}

func (b *BangumiTorrentInfo) SetENGName(name string) {
	b.detail.TorrentEngName = name
}

func (b *BangumiTorrentInfo) SetJPNName(name string) {
	b.detail.TorrentJpnName = name
}

func (b *BangumiTorrentInfo) GetContent() string {
	return util.GetJsonStrFromStruct(b.content)
}

func (b *BangumiTorrentInfo) GetDetail() string {
	return util.GetJsonStrFromStruct(b.detail)
}
