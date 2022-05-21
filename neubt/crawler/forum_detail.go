package crawler

import (
	"crawler/neubt"
	"crawler/neubt/html"
	"log"
	"strconv"
)

type ForumDetail interface {
	// GetFloorDetailFromForum return all the floors in a forum
	GetFloorDetailFromForum(link string) ([]*FloorDetail, error)
	// GetTorrentURLFromForum get all floors that contain a torrent-like link
	GetTorrentURLFromForum(link string) ([]*FloorDetail, error)
}

type ForumDetailImpl struct {
	client   neubt.Client
	MaxDepth int
}

func NewForumDetail(client neubt.Client) ForumDetail {
	f := &ForumDetailImpl{
		client:   client.Clone(),
		MaxDepth: 10,
	}
	return f
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

func (f *ForumDetailImpl) GetFloorDetailFromForum(link string) ([]*FloorDetail, error) {
	var floors []*FloorDetail
	nextPage := link
	for i := 0; i < f.MaxDepth && nextPage != ""; i++ {
		resp, err := f.client.SyncVisit(nextPage)
		if err != nil {
			return nil, err
		}
		node, err := html.NewNodeFromBytes(resp.Body)
		if err != nil {
			return nil, err
		}
		nodeList, err := node.GetInnerNodeList(`//div[@id="postlist"]/div[starts-with(@id, 'post_')]`)
		if err != nil {
			return nil, err
		}
		for _, n := range nodeList {
			floor := &FloorDetail{}
			floor.Comment, err = f.getCommentDetailByDiv(n)
			floor.UserInfo, err = f.getUserDetailByDiv(n)
			if err != nil {
				log.Println("error in GetDetailFrom", err)
				continue
			}
			floors = append(floors, floor)
		}
		nextPage, err = node.GetInnerString(`.//*[text()="下一页"]//@href`)
		// no next page
		if err != nil {
			break
		}
	}
	return floors, nil
}

func (f *ForumDetailImpl) getUserDetailByDiv(n *html.NodeDecorator) (*UserInfo, error) {
	node, err := n.GetInnerNode(`.//div[@class="pls favatar"]`)
	if err != nil {
		return nil, err
	}
	userInfo := &UserInfo{}
	userInfo.Username, _ = node.GetInnerString(`.//a[@class="xw1"]`)
	userInfo.Homepage, _ = node.GetInnerString(`.//a[@class="xw1"]/@href`)
	pointStr, _ := node.GetInnerString(`.//*[text()="积分"]//@title`)
	userInfo.Point, _ = strconv.Atoi(pointStr)
	return userInfo, nil
}

func (f *ForumDetailImpl) getCommentDetailByDiv(n *html.NodeDecorator) (*CommentInfo, error) {
	commentInfo := &CommentInfo{}
	commentInfo.Text, _ = n.GetInnerString(`.//td[starts-with(@id, 'postmessage_')]`)
	// get torrent file
	commentInfo.TorrentName, _ = n.GetInnerString(`.//*[@class="attnm"]/a`)
	commentInfo.TorrentLink, _ = n.GetInnerString(`.//*[@class="attnm"]/a/@href`)
	return commentInfo, nil
}

// GetTorrentURLFromForum get all floor that contains a torrent-like link
func (f *ForumDetailImpl) GetTorrentURLFromForum(link string) ([]*FloorDetail, error) {
	floors, err := f.GetFloorDetailFromForum(link)
	if err != nil {
		return nil, err
	}
	var torrents []*FloorDetail
	for _, floor := range floors {
		if floor.Comment.TorrentLink != "" {
			torrents = append(torrents, floor)
		}
	}
	return torrents, nil
}
