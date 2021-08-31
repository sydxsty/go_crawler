package main

import (
	"goCrawler/controller"
	"log"
)

func main() {
	if err := controller.Login(); err != nil {
		log.Fatal(err)
	}
	controller.CrawlBangumiInfo()
}
