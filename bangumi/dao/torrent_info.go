package dao

import (
	bgm "crawler/bangumi/anime_control"
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
	Content   []interface{}
	Detail    *BangumiTorrentDetail
}

type BangumiTorrentDetail struct {
	torrentDownloadURL string
	torrentChsName     string
	torrentEngName     string
	torrentJpnName     string
	teamName           string
	resolution         string
	format             string
	language           string
	episode            string
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
	if b.Content, ok = raw["content"].([]interface{}); !ok {
		return nil, errors.New("cannot get torrent content")
	}
	return b, nil
}

func (b *BangumiTorrentInfo) InitTorrentDetail(miscList []map[string]interface{}) {
	d := &BangumiTorrentDetail{}
	b.Detail = d
	d.torrentDownloadURL = `/download/torrent/` + b.TorrentId + `/` + b.TorrentId + `.torrent`
	for _, misc := range miscList {
		for _, cate := range b.TagIds {
			if misc["_id"] == cate {
				switch misc["type"] {
				case "bangumi":
					translate := misc["locale"].(map[string]interface{})
					d.torrentChsName, _ = translate["zh_cn"].(string)
					d.torrentEngName, _ = translate["en"].(string)
					d.torrentJpnName, _ = translate["ja"].(string)
				case "lang":
					if len(d.language) != 0 {
						d.language += " "
					}
					d.language += misc["name"].(string)
				default:
					break
				}
			}
		}
	}
	// TODO: get language detail
	// TODO: get bangumi anime name detail
	d.resolution = getString(b.bgmFilter.GetResolution(b.Title))
	d.format = getString(b.bgmFilter.GetMediaInfo(b.Title))
	d.teamName = getString(b.bgmFilter.GetTeam(b.Title))
	d.episode = getString([]string{
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
	if b.Detail == nil {
		return "", errors.New("not init")
	}
	if len(b.Detail.torrentChsName) != 0 {
		return b.Detail.torrentDownloadURL, nil
	}
	return "", errors.New("can not get torrent url")
}

func (b *BangumiTorrentInfo) MustGetCHSName() string {
	if b.Detail == nil {
		return ""
	}
	if len(b.Detail.torrentChsName) != 0 {
		return b.Detail.torrentChsName
	}
	if len(b.Detail.torrentEngName) != 0 {
		return b.Detail.torrentEngName
	}
	return ""
}

func (b *BangumiTorrentInfo) MustGetENGName() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.torrentEngName
}

func (b *BangumiTorrentInfo) MustGetJPNName() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.torrentJpnName
}

func (b *BangumiTorrentInfo) MustGetTeamName() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.teamName
}

func (b *BangumiTorrentInfo) MustGetResolution() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.resolution
}

func (b *BangumiTorrentInfo) MustGetEpisode() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.episode
}

func (b *BangumiTorrentInfo) MustGetFormat() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.format
}

func (b *BangumiTorrentInfo) MustGetLanguage() string {
	if b.Detail == nil {
		return ""
	}
	return b.Detail.language
}

func (b *BangumiTorrentInfo) UpdateFinalInformation(overrideCHSName func() (string, error)) {
	if overrideCHSName != nil {
		name, err := overrideCHSName()
		if err == nil && name != "" {
			log.Println("set torrent chsName to: ", name)
			b.Detail.torrentChsName = name
		}
	}
	target := b.bgmFilter.GetSeasonType(b.Detail.torrentChsName)
	if len(target) == 0 { // the season is empty
		b.Detail.torrentChsName += getString(b.bgmFilter.GetSeasonType(b.Title))
	}
}
