package dao

import (
	bgm "crawler/bangumi/anime_control"
	"crawler/util"
	"github.com/pkg/errors"
	"log"
)

type BangumiTorrentInfo struct {
	bgmFilter *bgm.BangumiFilter
	Title     string // raw title
	TorrentId string
	TeamId    string
	TagIds    []interface{}
	InfoHash  string
	content   *TorrentFileList
	detail    *BangumiTorrentDetail
}

type BangumiTorrentDetail struct {
	DownloadURL string
	ChsName     string
	EngName     string
	JpnName     string
	Teams       []string
	Resolution  []string
	Format      []string
	Languages   []string
	SEInfo      *SEInfo
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
	tc, ok := raw["content"].([]interface{})
	if !ok {
		return nil, errors.New("cannot get torrent content")
	}
	var err error
	b.content, err = NewTorrentFileList(tc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to covert torrent content")
	}
	return b, nil
}

func (b *BangumiTorrentInfo) InitTorrentDetail(miscList []map[string]interface{}) {
	d := &BangumiTorrentDetail{}
	b.detail = d
	d.DownloadURL = `/download/torrent/` + b.TorrentId + `/` + b.TorrentId + `.torrent`
	for _, misc := range miscList {
		for _, cate := range b.TagIds {
			if misc["_id"] == cate {
				switch misc["type"] {
				case "bangumi":
					translate := misc["locale"].(map[string]interface{})
					d.ChsName, _ = translate["zh_cn"].(string)
					d.EngName, _ = translate["en"].(string)
					d.JpnName, _ = translate["ja"].(string)
				case "lang":
					d.Languages = append(d.Languages, misc["name"].(string))
				default:
					break
				}
			}
		}
	}
	// TODO: get language detail
	d.Resolution = b.bgmFilter.GetResolution(b.Title)
	d.Format = b.bgmFilter.GetMediaInfo(b.Title)
	d.Teams = b.bgmFilter.GetTeam(b.Title)
	d.SEInfo = NewSEInfoFromTitle(b.Title, b.bgmFilter)
}

func getString(strList []string) string {
	return getStringWithDelim(strList, " ")
}

func getStringWithDelim(strList []string, delim string) string {
	res := ""
	for _, str := range strList {
		if str == "" {
			continue
		}
		if len(res) != 0 {
			res += delim
		}
		res += str
	}
	return res
}

func (b *BangumiTorrentInfo) GetTorrentDownloadURL() (string, error) {
	if b.detail == nil {
		return "", errors.New("not init")
	}
	if len(b.detail.ChsName) != 0 {
		return b.detail.DownloadURL, nil
	}
	return "", errors.New("can not get torrent url")
}

func (b *BangumiTorrentInfo) MustGetCHSName() string {
	if b.detail == nil {
		return ""
	}
	if len(b.detail.ChsName) != 0 {
		return b.detail.ChsName
	}
	if len(b.detail.EngName) != 0 {
		return b.detail.EngName
	}
	return ""
}

func (b *BangumiTorrentInfo) MustGetENGName() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.EngName
}

func (b *BangumiTorrentInfo) MustGetJPNName() string {
	if b.detail == nil {
		return ""
	}
	return b.detail.JpnName
}

func (b *BangumiTorrentInfo) MustGetTeamStr() string {
	if b.detail == nil {
		return ""
	}
	return getStringWithDelim(b.detail.Teams, "&")
}

func (b *BangumiTorrentInfo) MustGetTeam() []string {
	if b.detail == nil {
		return nil
	}
	return b.detail.Teams
}

func (b *BangumiTorrentInfo) MustGetResolution() string {
	if b.detail == nil {
		return ""
	}
	return getString(b.detail.Resolution)
}

func (b *BangumiTorrentInfo) MustGetEpisode() string {
	if b.detail == nil {
		return ""
	}
	return getString(b.detail.SEInfo.GetEpisodeStringList())
}

func (b *BangumiTorrentInfo) MustGetFormat() string {
	if b.detail == nil {
		return ""
	}
	return getString(b.detail.Format)
}

func (b *BangumiTorrentInfo) MustGetLanguage() string {
	if b.detail == nil {
		return ""
	}
	return getString(b.detail.Languages)
}

// SetReleaseCHSName change the chs name to release version, the chs name may not be searchable in ptgen
func (b *BangumiTorrentInfo) SetReleaseCHSName(name string) {
	if len(name) != 0 {
		log.Println("set torrent chsName to: ", name)
		b.detail.ChsName = name
	}
	target := b.bgmFilter.GetSeasonType(b.detail.ChsName)
	if season := getString(b.bgmFilter.GetSeasonType(b.Title)); len(target) == 0 && len(season) != 0 {
		b.detail.ChsName += " " + season
	}
}

func (b *BangumiTorrentInfo) SetCHSName(name string) {
	b.detail.ChsName = name
}

func (b *BangumiTorrentInfo) SetENGName(name string) {
	b.detail.EngName = name
}

func (b *BangumiTorrentInfo) SetJPNName(name string) {
	b.detail.JpnName = name
}

func (b *BangumiTorrentInfo) GetContent() string {
	str, err := b.content.PrintToString(10)
	if err != nil {
		log.Println(err)
		return ""
	}
	return str
}

func (b *BangumiTorrentInfo) GetDetail() string {
	return util.GetJsonStrFromStruct(b.detail)
}

func (b *BangumiTorrentInfo) ContainsFinishedSeasons() bool {
	if b.detail == nil {
		return false
	}
	return b.detail.SEInfo.Finished
}

func (b *BangumiTorrentInfo) ContainsMovie() bool {
	if b.detail == nil {
		return false
	}
	if len(b.detail.SEInfo.Movie) > 0 {
		return true
	}
	return false
}
