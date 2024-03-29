package ptgen

import (
	"crawler/ptgen/bgmtv"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type PTGen interface {
	GetBangumiLinkByNames(jpnName string, names ...string) ([]*BangumiLinkDetail, error)
	GetBangumiLinkByName(name string) ([]*BangumiLinkDetail, error)
	GetBangumiInfoByLink(link string) (map[string]interface{}, error)
}

type BangumiLinkDetail struct {
	ChnName string
	JpnName string
	Link    string
}

type BangumiInfoDetail struct {
	ChnName string
	EngName string
	JpnName string
	Detail  string
}

type PTGenImpl struct {
	client    Client
	tvClient  bgmtv.Client
	apiClient bgmtv.Client
}

func NewPTGen() (PTGen, error) {
	apiClient, err := bgmtv.NewAPIClient()
	if err != nil {
		return nil, err
	}
	tvClient, err := bgmtv.NewClient()
	if err != nil {
		return nil, err
	}
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &PTGenImpl{
		client:    client,
		tvClient:  tvClient,
		apiClient: apiClient,
	}, nil
}

func (p *PTGenImpl) GetBangumiLinkByNames(jpnName string, names ...string) ([]*BangumiLinkDetail, error) {
	for _, name := range append([]string{jpnName}, names...) {
		result, err := p.GetBangumiLinkByName(name)
		if err == nil {
			return result, nil
		}
		log.Println("this name contains no result, switch to another name, ", name, err)
	}
	return nil, errors.New("can not get any valid result")
}

func (p *PTGenImpl) searchBangumi(name string) ([]interface{}, error) {
	resp, err := p.client.SyncVisit(`/?` + `search=` + url.QueryEscape(name) + `&source=bangumi`)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}
	msg, ok := result["error"]
	if !ok || msg != nil {
		return nil, errors.New("remote server failure: " + msg.(string))
	}
	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, errors.New("can not unpack result")
	}
	return data, nil
}

func (p *PTGenImpl) GetBangumiLinkByName(name string) ([]*BangumiLinkDetail, error) {
	if name == "" {
		return nil, errors.New("query name is empty")
	}
	name = strings.ReplaceAll(name, "!", "！")
	name = strings.ReplaceAll(name, "'", "’")
	log.Println("using local bgmtv client, ", name)
	data, err := bgmtv.SearchBangumi(p.apiClient, name) // local first
	if err != nil {
		log.Println("switch to cf worker client, ", name)
		data, err = p.searchBangumi(name)
		if err != nil {
			return nil, errors.New("can not GetBangumiLinkByName from local or cf")
		}
	}
	links := make([]*BangumiLinkDetail, 0)
	for _, node := range data {
		unmarshalNode := node.(map[string]interface{})
		if unmarshalNode["subtype"].(string) == "动画/二次元番" {
			links = append(links, &BangumiLinkDetail{
				ChnName: unmarshalNode["title"].(string),
				JpnName: unmarshalNode["subtitle"].(string),
				Link:    unmarshalNode["link"].(string),
			})
			break
		}
	}
	if links == nil || len(links) == 0 {
		return nil, errors.New("result is empty")
	}
	return links, nil
}

func genBangumi(client Client, link string) (map[string]interface{}, error) {
	resp, err := client.SyncVisit(`/?url=` + link)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PTGenImpl) GetBangumiInfoByLink(link string) (map[string]interface{}, error) {
	log.Println("using local bgmtv client, ", link)
	result, err := bgmtv.GenBangumi(p.apiClient, link) // TODO: use tv client
	if err != nil {
		log.Println("switch to cf worker client, ", link)
		result, err = genBangumi(p.client, link)
		if err != nil {
			return nil, errors.New("can not GetBangumiInfoByLink from local or cf")
		}
	}
	success, ok := result["success"].(bool)
	if !ok || !success {
		return nil, errors.New("cloudflare worker network error")
	}
	return result, nil
}

func GetDetailFromInfo(info map[string]interface{}) (*BangumiInfoDetail, error) {
	b := &BangumiInfoDetail{}
	b.Detail = getFormatFromInfo(info)
	if len(b.Detail) == 0 {
		return nil, errors.New("covert failure")
	}
	vs, ok := info["info"].([]interface{})
	if !ok {
		return b, nil
	}
	b.EngName = getENGNameFromInfo(vs)
	b.ChnName = getCHNNameFromInfo(vs)
	b.JpnName = getJPNNameFromInfo(vs)
	return b, nil
}

func getFormatFromInfo(info map[string]interface{}) string {
	format, ok := info["format"].(string)
	if !ok {
		return ""
	}
	return strings.ReplaceAll(format, " ", "")
}

func getENGNameFromInfo(info []interface{}) string {
	nonEngFilter := regexp.MustCompile(`[\x{4e00}-\x{9fa5}]|[\x{3040}-\x{309F}\x{30A0}-\x{30FF}\x{31F0}-\x{31FF}]`)
	for _, v := range info {
		str := v.(string)
		if strings.Contains(str, "别名: ") {
			str = strings.ReplaceAll(str, "别名: ", "")
			if len(nonEngFilter.FindAllString(str, -1)) == 0 {
				// non chinese or japanese
				return str
			}
		}
	}
	return ""
}

func getJPNNameFromInfo(info []interface{}) string {
	jpnFilter := regexp.MustCompile(`[\x{3040}-\x{309F}\x{30A0}-\x{30FF}\x{31F0}-\x{31FF}]`)
	for _, v := range info {
		str := v.(string)
		if strings.Contains(str, "别名: ") {
			str = strings.ReplaceAll(str, "别名: ", "")
			if len(jpnFilter.FindAllString(str, -1)) != 0 {
				// japanese
				return str
			}
		}
	}
	return ""
}

func getCHNNameFromInfo(info []interface{}) string {
	for _, v := range info {
		str := v.(string)
		if strings.Contains(str, "中文名: ") {
			return strings.ReplaceAll(str, "中文名: ", "")
		}
	}
	return ""
}
