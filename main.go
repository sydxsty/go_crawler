package main

import (
	"goCrawler/controller"
	"log"
	"time"
)

func main() {
	if err := controller.Login(); err != nil {
		log.Fatal(err)
	}
	for {
		controller.CrawlBangumiInfo()
		log.Println("sleep 600 sec to continue")
		time.Sleep(time.Second * 600)
	}
}
