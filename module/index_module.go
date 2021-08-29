package module

import (
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"log"
	"regexp"
)

type IndexModule interface {
	ScraperModule
	GetResourceIndex() []*dao.TorrentInfo
	GetForum(url string) []*dao.TorrentInfo
	GetForumList() map[string]string
}

type indexModuleImpl struct {
	scraperModuleImpl
}

func NewIndexModule() IndexModule {
	c := &indexModuleImpl{}
	c.init()
	return c
}

func (c *indexModuleImpl) GetResourceIndex() []*dao.TorrentInfo {
	collector := c.getClonedCollector()
	// analyse resource index
	var torrentInfoList []*dao.TorrentInfo
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		threadNodes := node.GetInnerNodeList(`//a[@class="s xst"]/../..`)
		// get all thread node
		for _, n := range threadNodes {
			torrentInfo := c.indexAnalysis(n)
			if torrentInfo == nil {
				continue
			}
			torrentInfoList = append(torrentInfoList, torrentInfo)
		}
	})
	if err := collector.Visit(c.getAbsoluteURL(`plugin.php?id=neubt_resourceindex`)); err != nil {
		log.Fatal(err)
	}
	return torrentInfoList
}

func (c *indexModuleImpl) GetForum(url string) []*dao.TorrentInfo {
	collector := c.getClonedCollector()
	// analyse resource index
	var torrentInfoList []*dao.TorrentInfo
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		threadNodes := node.GetInnerNodeList(`//*[starts-with(@id, 'normalthread_')]//a[@class="s xst"]/../..`)
		for _, n := range threadNodes {
			torrentInfo := c.forumAnalysis(n)
			if torrentInfo == nil {
				continue
			}
			torrentInfoList = append(torrentInfoList, torrentInfo)
		}
	})
	if err := collector.Visit(c.getAbsoluteURL(url)); err != nil {
		log.Fatal(err)
	}
	return torrentInfoList
}

func (c *indexModuleImpl) GetForumList() map[string]string {
	collector := c.getClonedCollector()
	forumList := make(map[string]string)
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// purge link
		if matchString, _ := regexp.MatchString(`forum-([0-9]*)-1.html`, link); len(e.Text) == 0 || !matchString {
			return
		}
		forumList[e.Text] = e.Request.AbsoluteURL(link)
	})
	if err := collector.Visit(c.getAbsoluteURL("/")); err != nil {
		log.Fatal(err)
	}
	return forumList
}

func (c *indexModuleImpl) indexAnalysis(node *NodeDecorator) *dao.TorrentInfo {
	raw := &dao.RawTorrentInfo{}
	// generate a detailed information of a torrent
	raw.Title = node.GetInnerNode(`.//a[@class="s xst"]`).GetString()
	raw.Link = node.GetInnerNode(`.//a[@class="s xst"]/@href`).GetString()
	raw.Size = node.GetInnerNode(`./td[3]`).GetString()
	raw.Discount = node.GetInnerNode(`.//a[@class="s xst"]/../img/@src`).GetString()
	raw.Signal = node.GetInnerNode(`.//strong/img/@src`).GetString()
	raw.Uploader = node.GetInnerNode(`./td[last()]/a`).GetString()
	raw.UploaderPage = node.GetInnerNode(`./td[last()]/a/@href`).GetString()
	return dao.NewAggregatedTorrentInfo(raw)
}

func (c *indexModuleImpl) forumAnalysis(node *NodeDecorator) *dao.TorrentInfo {
	raw := &dao.RawTorrentInfo{}
	// generate a detailed information of a torrent
	raw.Title = node.GetInnerNode(`.//a[@class="s xst"]`).GetString()
	raw.Link = node.GetInnerNode(`.//a[@class="s xst"]/@href`).GetString()
	raw.Size = node.GetInnerNode(`./td[3]`).GetString()
	raw.Discount = node.GetInnerNode(`.//a[@class="s xst"]/../img[last()]/@src`).GetString()
	raw.Signal = node.GetInnerNode(`./td[2]/img/@src`).GetString()
	raw.Uploader = node.GetInnerNode(`./td[last()]//a`).GetString()
	raw.UploaderPage = node.GetInnerNode(`./td[last()]//a/@href`).GetString()
	return dao.NewAggregatedTorrentInfo(raw)
}
