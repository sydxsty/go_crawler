package main

import (
	"crawler/neubt"
	"crawler/neubt/crawler"
	"crawler/neubt/dao"
	"crawler/util/html"
	"log"
	"os"
	"time"
)

func main() {
	r := NewRSS()
	for {
		log.Println("crawling film(forum-45-1) and resource index")
		r.crawlForum("forum-45-1.html")
		r.crawlForum("forum-13-1.html")
		r.crawlResourceIndex()
		log.Println("sleep 600 sec to continue")
		time.Sleep(time.Second * 600)
	}
}

type RSS struct {
	*neubt.NeuBT
}

func NewRSS() *RSS {
	return &RSS{
		neubt.NewNeuBT(),
	}
}

func (r *RSS) crawlForum(link string) {
	forum := crawler.NewForum(r.Client)
	list, err := forum.GetForum(link)
	if err != nil {
		log.Fatal(err, "can not get forum")
	}
	r.downloadTorrentListFromNodes(list)
}

func (r *RSS) crawlResourceIndex() {
	ri := crawler.NewResourceIndex(r.Client)
	list, err := ri.GetResourceIndex()
	if err != nil {
		log.Fatal(err, "can not get resource index")
	}
	r.downloadTorrentListFromNodes(list)
}

func (r *RSS) downloadTorrentListFromNodes(nodes []*html.NodeDecorator) {
	results, err := dao.NodeListToTorrentInfoList(r.KVS, nodes)
	if err != nil {
		log.Fatal(err, "can not cover NodeListToTorrentInfoList")
	}
	for _, result := range results {
		r.downloadTorrentByInfo(result)
	}
}

func (r *RSS) downloadTorrentByInfo(info *dao.TorrentInfo) {
	/* please apply your custom filter here */
	if info.Crawled || info.TorrentID < r.Config.ThreadWaterMark || info.Discount < r.Config.DiscountWaterMark {
		log.Println("torrent condition not met")
		return
	}
	detail := crawler.NewThreadDetail(r.Client)
	torrentURLs, err := detail.GetFloorDetailFromThread(info.Link)
	if err != nil {
		log.Println(err, "wrong torrent info")
	}
	downloader := crawler.NewDownloader(r.Client)
	for _, torrentURL := range torrentURLs {
		data, err := downloader.DownloadFromNestedURL(torrentURL.Comment.TorrentLink)
		if err != nil {
			log.Println(err, "can not download torrent from url")
			continue
		}
		err = os.WriteFile(r.Config.TorrentPath+torrentURL.Comment.TorrentName, data, 0666)
		if err != nil {
			log.Println(err, "write file failure")
			continue
		}
		err = r.Webui.AddTorrentFromFile(r.Config.TorrentPath, torrentURL.Comment.TorrentName)
		if err != nil {
			log.Println(err, "failed to upload torrent to webui")
			continue
		}
		info.Crawled = true
		log.Println("downloaded: ", info)
		err = r.KVS.Put(info.Link, info)
		if err != nil {
			log.Println("can not update torrent info, ", err)
		}
	}
}
