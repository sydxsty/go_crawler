package ptgen

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

type PTGen interface {
	GetBangumiLinkByNames(jpnName string, names ...string) (map[string]string, error)
	GetBangumiLinkByName(name string) (map[string]string, error)
	GetBangumiDetailByLink(link string) (map[string]interface{}, error)
}

type PTGenImpl struct {
	client Client
}

func NewPTGen(client Client) PTGen {
	return &PTGenImpl{
		client: client.Clone(),
	}
}

func (p *PTGenImpl) GetBangumiLinkByNames(jpnName string, names ...string) (map[string]string, error) {
	for _, name := range append([]string{jpnName}, names...) {
		result, err := p.GetBangumiLinkByName(name)
		if err == nil {
			return result, nil
		}
		log.Println("this name contains no result, switch to another name, ", name, err)
	}
	return nil, errors.New("can not get any valid result")
}

func (p *PTGenImpl) GetBangumiLinkByName(name string) (map[string]string, error) {
	if name == "" {
		return nil, errors.New("query name is empty")
	}
	resp, err := p.client.SyncVisit(`/?` + `search=` + url.QueryEscape(name) + `&source=bangumi`)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, errors.New("can not unpack result")
	}
	linkMap := make(map[string]string)
	for _, node := range data {
		unmarshalNode := node.(map[string]interface{})
		if unmarshalNode["subtype"].(string) == "动画/二次元番" {
			linkMap[unmarshalNode["title"].(string)] = unmarshalNode["link"].(string)
			break
		}
	}
	if linkMap == nil || len(linkMap) == 0 {
		return nil, errors.New("result is empty")
	}
	return linkMap, nil
}

func (p *PTGenImpl) GetBangumiDetailByLink(link string) (map[string]interface{}, error) {
	resp, err := p.client.SyncVisit(`/?url=` + link)
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

func GetTextFromDetail(detail map[string]interface{}) (string, error) {
	value, ok := detail["format"].(string)
	if !ok {
		return "", errors.New("covert failure")
	}
	value = strings.ReplaceAll(value, " ", "")
	return value, nil
}
