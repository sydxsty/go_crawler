package main

import (
	"goCrawler/dao"
	"goCrawler/module"
	"log"
)

func downloadInfoList(infoList []*dao.TorrentInfo) {
	detail := module.NewDetailModule()
	download := module.NewDownloader()
	for _, info := range infoList {
		if info.Crawled || info.TorrentID < dao.YAMLConfig.ThreadWaterMark || info.Discount < dao.YAMLConfig.DiscountWaterMark {
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
				log.Println("downloading: ")
				log.Println(info)
				if err := dao.SaveTorrentInfoToDB(info); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

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
	downloadInfoList(c.GetForum("forum-45-1.html"))
	downloadInfoList(c.GetResourceIndex())
}
