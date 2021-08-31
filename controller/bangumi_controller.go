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
	result := SearchWithPTGen(info)
	ptGen := module.NewPTGen()
	for _, v := range result {
		return ptGen.GetBangumiDetailByLink(v)
	}
	return nil
}

func SearchWithPTGen(info *dao.BangumiTorrentInfo) map[string]string {
	var result map[string]string
	if len(info.Detail.TorrentChsName) != 0 {
		if err := dao.LoadFromDB("ptgen_error_name_"+info.Detail.TorrentChsName, &result); err == nil {
			return nil
		}
		if err := dao.LoadFromDB("ptgen_info_name_"+info.Detail.TorrentChsName, &result); err == nil {
			return result
		}
	}
	ptGen := module.NewPTGen()
	if result == nil {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentJpnName)
	}
	if result == nil {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentChsName)
	}
	if result == nil {
		result = ptGen.GetBangumiLinkByName(info.Detail.TorrentEngName)
	}
	if result == nil || len(result) == 0 {
		log.Println("torrent name is empty")
		dao.SaveToDB("ptgen_error_name_"+info.Detail.TorrentChsName, result)
		return nil
	} else {
		if len(info.Detail.TorrentChsName) != 0 {
			dao.SaveToDB("ptgen_info_name_"+info.Detail.TorrentChsName, result)
		}
		return result
	}
}
