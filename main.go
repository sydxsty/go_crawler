package main

import (
	"goCrawler/controller"
	"goCrawler/dao"
	"goCrawler/module"
	"log"
	"time"
)

func main() {
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
	for {
		log.Println("crawling film(forum-45-1) and resource index")
		controller.DownloadTorrentByInfoList(c.GetForum("forum-45-1.html"))
		controller.DownloadTorrentByInfoList(c.GetResourceIndex())
		log.Println("sleep 600 sec to continue")
		time.Sleep(time.Second * 600)
	}
}

func bangumiCrawler() {
	if err := controller.Login(); err != nil {
		log.Fatal(err)
	}
	for {
		controller.CrawlBangumiInfo()
		log.Println("sleep 600 sec to continue")
		time.Sleep(time.Second * 600)
	}
}
