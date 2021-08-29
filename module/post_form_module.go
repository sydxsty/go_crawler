package module

import (
	"github.com/gocolly/colly/v2"
	"log"
)

type ForumModule interface {
	ScraperModule
}

type forumModuleImpl struct {
	scraperModuleImpl
	tidList map[string]string

	formHash string
	postTime string
	wysiwyg string
	special string
	specialExtra string
	tid string
	subject string

	torrentFileName string
	torrentFileBytes []byte

	text string
}

func NewForumModule(fid string) ForumModule {
	c := &forumModuleImpl{}
	c.init()
	collector := c.getClonedCollector()
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.formHash = node.GetInnerNode(`.//input[@id="formhash"]/@value`).GetString()
		c.postTime = node.GetInnerNode(`.//input[@id="posttime"]/@value`).GetString()
		c.wysiwyg = node.GetInnerNode(`.//input[@name="wysiwyg"]/@value`).GetString()
		c.special = node.GetInnerNode(`.//input[@name="special"]/@value`).GetString()
		c.specialExtra = node.GetInnerNode(`.//input[@name="specialextra"]/@value`).GetString()
		tidNodeList := node.GetInnerNodeList(`.//select[@name="typeid"]/option`)
		c.tidList = make(map[string]string)
		for _, tidNode := range tidNodeList {
			c.tidList[tidNode.GetString()] = tidNode.GetInnerNode(`./@value`).GetString()
		}

		if v, ok := c.tidList["选择主题分类"]; !ok {
			log.Fatal("no matching default tid")
		} else {
			c.tid = v
		}
		// default subject name
		c.subject =  node.GetInnerNode(`.//input[@name="subject"]/@value`).GetString()

	})
	url := `forum.php?mod=post&action=newthread&fid=` + fid + `&specialextra=torrent`
	if err := collector.Visit(c.getAbsoluteURL(url)); err != nil {
		log.Fatal(err)
	}
	return c
}