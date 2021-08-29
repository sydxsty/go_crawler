package controller

import (
	"goCrawler/dao"
	"goCrawler/module"
	"log"
)

func GetListByForumName(nameList ...string) []*dao.TorrentInfo {
	var infoList []*dao.TorrentInfo
	c := module.NewIndexModule()
	forumMap := c.GetForumList()
	log.Println("all available forums:")
	log.Println(forumMap)
	for _, name := range nameList {
		if url, ok := forumMap[name]; ok {
			infoList = append(infoList, c.GetForum(url)...)
		}
	}
	return infoList
}
