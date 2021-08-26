package main

import (
	"goCrawler/dao"
	"goCrawler/module"
	"log"
)

func main() {
	// test case
	c := module.NewIndexModule()
	if dao.YAMLConfig.UseCookie {
		if err := c.LoadCookie(); err != nil {
			log.Fatal(err)
		}
	} else {
		c.Login(dao.YAMLConfig.Username, dao.YAMLConfig.Password)
		if err := c.SaveCookie(); err != nil {
			log.Fatal(err)
		}
	}
	//infoList := c.GetResourceIndex()
	infoList := c.GetForum("forum-45-1.html")
	detail := module.NewDetailModule()
	download := module.NewDownloader()
	for _, info := range infoList {
		if info.Crawled {
			continue
		}
		form := detail.GetDetailFrom(info)
		for _, floor := range form.Floors {
			if floor.Comment.TorrentLink != "" {
				if err := download.DownloadTorrent(
					floor.Comment.TorrentLink,
					floor.Comment.TorrentName,
					dao.YAMLConfig.TorrentPath); err != nil {
					log.Println(err)
					continue
				}
				if err := download.AddTorrentToBitTorrent(
					floor.Comment.TorrentName,
					dao.YAMLConfig.TorrentPath); err != nil {
					continue
				}
				info.Crawled = true
				if err := dao.SaveTorrentInfoToDB(info); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
