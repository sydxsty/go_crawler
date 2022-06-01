package bgmtv

import (
	"crawler/util/html"
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"regexp"
	"strings"
)

func SearchBangumi(query string) ([]interface{}, error) {
	tpDict := map[int]string{1: "漫画/小说", 2: "动画/二次元番", 3: "音乐", 4: "游戏", 6: "三次元番"}
	client, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	r, err := getJsonSearchResp(client, `search/subject/`+url.QueryEscape(query)+`?responseGroup=large`, 3)
	if err != nil {
		return nil, errors.Wrap(err, "can not get response")
	}
	rl, ok := r["list"].([]interface{})
	if !ok || len(rl) == 0 {
		return nil, errors.New("result list is empty")
	}
	var resultList []interface{}
	for _, v := range rl {
		result := make(map[string]interface{})
		uv := v.(map[string]interface{})
		result["year"], _ = uv["air_date"]
		result["subtype"], _ = tpDict[int(uv["type"].(float64))]
		if uv["name_cn"] != "" {
			result["title"] = uv["name_cn"]
		} else {
			result["title"] = uv["name"]
		}
		result["subtitle"] = uv["name"]
		result["link"] = uv["url"]
		resultList = append(resultList, result)
	}
	return resultList, nil
}

func getJsonSearchResp(client Client, link string, retry uint) (map[string]interface{}, error) {
	for i := 0; uint(i) < retry; i++ {
		resp, err := client.SyncVisit(link)
		if err != nil {
			return nil, errors.Wrap(err, "can not get response")
		}
		var r map[string]interface{}
		err = json.Unmarshal(resp.Body, &r)
		if err != nil {
			continue
		}
		return r, nil
	}
	return nil, errors.Errorf("can not get valid resp within %d times.", retry)
}

func GenBangumi(client Client, link string) (interface{}, error) {
	resp, err := client.SyncVisit(link)
	if err != nil {
		return nil, errors.Wrap(err, "can not load page")
	}
	notExistFilter := regexp.MustCompile(`呜咕，出错了`)
	if notExistFilter.FindAllString(string(resp.Body), -1) != nil {
		return nil, errors.New("404 not exist")
	}
	data := make(map[string]interface{})
	data["site"] = "bangumi"
	data["alt"] = link
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can not phrase response")
	}
	// 1. get all stuff
	coverStaffAnother, err := node.GetInnerNode(`.//div[@id="bangumiInfo"]`)
	if err != nil {
		return nil, errors.Wrap(err, "coverStaffAnother is empty")
	}
	// 1.1 get cover
	coverAnother, err := coverStaffAnother.GetInnerString(`.//a[@class="thickbox cover"]/@href`)
	if err != nil {
		return nil, errors.Wrap(err, "coverAnother is empty")
	}
	coverAnother = "https:" + regexp.MustCompile(`/cover/[lcmsg]/`).ReplaceAllString(coverAnother, "/cover/l/")
	data["cover"] = coverAnother
	data["poster"] = coverAnother
	// 1.2 get info
	infoAnother, err := coverStaffAnother.GetInnerNodeList(`.//ul[@id="infobox"]/li`)
	if err != nil {
		return nil, errors.Wrap(err, "infoAnother is empty")
	}
	staffFilter := regexp.MustCompile(`^(中文名|话数|放送开始|放送星期|别名|官方网站|播放电视台|其他电视台|Copyright)`)
	staff := make([]string, 0)
	info := make([]string, 0)
	for _, v := range infoAnother {
		str := v.GetString()
		if staffFilter.FindAllString(str, -1) == nil {
			staff = append(staff, str)
		} else {
			info = append(info, str)
		}
	}
	data["staff"] = staff
	data["info"] = info
	// 1.3 get story
	story, err := node.GetInnerString(`.//div[@id="subject_summary"]`)
	if err != nil {
		return nil, errors.Wrap(err, "storyAnother is empty")
	}
	data["story"] = story
	// 1.4 other stuff
	// TODO: bangumi_votes, bangumi_rating_average, tags
	// 2 get characters
	resp, err = client.SyncVisit(link + "/characters")
	if err != nil {
		return nil, err
	}
	node, err = html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can not get character response")
	}
	cast := make([]string, 0)
	castActors, err := node.GetInnerNodeList(`.//div[@id="columnInSubjectA"]/div[@class="light_odd"]/div[@class="clearit"]`)
	for _, v := range castActors {
		char, _ := v.GetInnerString(`.//h2/span[@class="tip"]`)
		cv, _ := v.GetInnerString(`.//div[@class="actorBadge clearit"]/p/small`)
		if len(cv) == 0 {
			cv, _ = v.GetInnerString(`.//div[@class="actorBadge clearit"]/p/a`)
		}
		if len(char) == 0 || len(cv) == 0 {
			continue
		}
		cast = append(cast, strings.ReplaceAll(char+`: `+cv, "/", ""))
	}
	data["cast"] = cast
	// 生成format
	var des string
	if len(coverAnother) != 0 {
		des += `[img]` + coverAnother + `[/img]\n\n`
	}
	if len(story) != 0 {
		des += `[b]Story: [/b]\n\n` + story + `\n\n`
	}
	if len(staff) != 0 {
		des += `[b]Staff: [/b]\n\n`
		for i, v := range staff {
			des += v + "\n"
			if i == 15 {
				break
			}
		}
		des += `\n\n`
	}

	if len(cast) != 0 {
		des += `[b]Cast: [/b]\n\n`
		for i, v := range cast {
			des += v + "\n"
			if i == 9 {
				break
			}
		}
		des += `\n\n`
	}
	des += `(来源于 ` + link + ` )\n`
	data["format"] = des
	data["success"] = true
	return data, nil
}
