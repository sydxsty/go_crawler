package module

import (
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"log"
	"strconv"
)

type DetailModule interface {
	ScraperModule
	GetDetailFrom(t *dao.TorrentInfo) *TorrentDetail
}

type detailModuleImpl struct {
	scraperModuleImpl
	MaxDepth int
}

type CommentInfo struct {
	Text        string
	TorrentName string
	TorrentLink string
}

type UserInfo struct {
	Username string
	Homepage string
	Point    int
}

type FloorDetail struct {
	Comment  *CommentInfo
	UserInfo *UserInfo
}

type TorrentDetail struct {
	Floors     []*FloorDetail
	ThreadInfo *dao.TorrentInfo
}

func NewDetailModule() DetailModule {
	d := &detailModuleImpl{}
	d.init()
	d.MaxDepth = 10
	return d
}

func (d *detailModuleImpl) GetDetailFrom(t *dao.TorrentInfo) *TorrentDetail {
	if t == nil {
		return nil
	}
	var floors []*FloorDetail
	nextPage := t.Link
	collector := d.getClonedCollector()
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		nodeList := node.GetInnerNodeList(`//div[@id="postlist"]/div[starts-with(@id, 'post_')]`)
		for _, n := range nodeList {
			// get user info
			floors = append(floors, &FloorDetail{
				Comment:  d.getCommentDetailByDiv(n),
				UserInfo: d.getUserDetailByDiv(n),
			})
		}
		nextPage = node.GetInnerNode(`.//*[text()="下一页"]//@href`).GetString()
	})
	for i := 0; i < d.MaxDepth && nextPage != ""; i++ {
		if err := collector.Visit(d.getAbsoluteURL(nextPage)); err != nil {
			break
		}
	}
	return &TorrentDetail{
		Floors:     floors,
		ThreadInfo: t,
	}
}

func (d *detailModuleImpl) getUserDetailByDiv(n *NodeDecorator) *UserInfo {
	node := n.GetInnerNode(`.//div[@class="pls favatar"]`)
	userInfo := &UserInfo{}
	userInfo.Username = node.GetInnerNode(`.//a[@class="xw1"]`).GetString()
	userInfo.Homepage = node.GetInnerNode(`.//a[@class="xw1"]/@href`).GetString()
	userInfo.Point, _ = strconv.Atoi(node.GetInnerNode(`.//*[text()="积分"]//@title`).GetString())
	return userInfo
}

func (d *detailModuleImpl) getCommentDetailByDiv(n *NodeDecorator) *CommentInfo {
	commentInfo := &CommentInfo{}
	commentInfo.Text = n.GetInnerNode(`.//td[starts-with(@id, 'postmessage_')]`).GetString()
	// get torrent file
	commentInfo.TorrentName = n.GetInnerNode(`.//*[@class="attnm"]/a`).GetString()
	commentInfo.TorrentLink = n.GetInnerNode(`.//*[@class="attnm"]/a/@href`).GetString()
	return commentInfo
}
