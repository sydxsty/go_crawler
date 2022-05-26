package crawler

import (
	"crawler/neubt"
	"crawler/neubt/html"
)

// ResourceIndex process pages like http://bt.neu6.edu.cn/plugin.php?id=neubt_resourceindex
type ResourceIndex interface {
	// GetResourceIndex return all node from / url
	GetResourceIndex() ([]*html.NodeDecorator, error)
}

type ResourceIndexImpl struct {
	client neubt.Client
}

func NewResourceIndex(client neubt.Client) ResourceIndex {
	r := &ResourceIndexImpl{
		client: client.Clone(),
	}
	return r
}

func (r ResourceIndexImpl) GetResourceIndex() ([]*html.NodeDecorator, error) {
	resp, err := r.client.SyncVisit(`plugin.php?id=neubt_resourceindex`)
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}
	return node.GetInnerNodeList(`//a[@class="s xst"]/../..`)
}
