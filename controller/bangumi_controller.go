package controller

import (
	"goCrawler/dao"
	"goCrawler/module"
	"log"
)

func GetAnimeNameList() []map[string]interface{} {
	t := module.NewBangumiModule()
	var tagList []string
	for _, id := range t.GetRecentTorrentList() {
		tagList = append(tagList, id.(map[string]interface{})["tag_id"].(string))
	}
	return t.GetAnimeMiscByTag(tagList...)
}

func GetLatestAnimeList() []*dao.BangumiTorrentInfo {
	t := module.NewBangumiModule()
	var bgmList []*dao.BangumiTorrentInfo
	for _, e := range t.GetLatestTorrentList() {
		bgm := dao.NewBangumiTorrentInfoFromMap(e.(map[string]interface{}))
		var tags []string
		for _, v := range bgm.TagIds {
			tags = append(tags, v.(string))
		}
		bgm.GenerateTorrentDetail(t.GetAnimeMiscByTag(tags...))
		bgmList = append(bgmList, bgm)
	}
	return bgmList
}

func GetTorrentPTGenDetail(info *dao.BangumiTorrentInfo) map[string]interface{} {
	ptGen := module.NewPTGen()
	var result map[string]string
	if len(info.Detail.TorrentJpnName) != 0 {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentJpnName)
	} else if len(info.Detail.TorrentChsName) != 0 {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentChsName)
	} else if len(info.Detail.TorrentEngName) != 0 {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentEngName)
	} else {
		log.Println("torrent name is empty")
	}
	for _, v := range result {
		ptResult := ptGen.GetBangumiDetailByLink(v)
		return ptResult
	}
	return nil
}
