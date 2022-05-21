package crawler

import (
	"crawler/neubt"
	"crawler/neubt/html"
)

type Forum interface {
	// GetForumList return all sections from / url
	GetForumList() (map[string]string, error)
	// GetForum return a list of torrent forums from specific sections
	GetForum(url string) ([]*html.NodeDecorator, error)
}

type ForumImpl struct {
	client neubt.Client
}

func NewForum(client neubt.Client) Forum {
	f := &ForumImpl{
		client: client.Clone(),
	}
	return f
}

func (f *ForumImpl) GetForumList() (map[string]string, error) {
	resp, err := f.client.SyncVisit("/")
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}
	nodes, err := node.GetInnerNodeList(`//dt/a`)
	if err != nil {
		return nil, err
	}
	forumList := make(map[string]string)
	for _, n := range nodes {
		key := n.GetString()
		value, _ := n.GetInnerString(`@href`)
		forumList[key] = value
	}
	return forumList, nil
}

func (f *ForumImpl) GetForum(url string) ([]*html.NodeDecorator, error) {
	resp, err := f.client.SyncVisit(url)
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}
	return node.GetInnerNodeList(`//*[starts-with(@id, 'normalthread_')]//a[@class="s xst"]/../..`)
}
