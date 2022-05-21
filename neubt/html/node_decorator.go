package html

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

type NodeDecorator struct {
	html.Node
}

func NewNodeFromBytes(data []byte) (*NodeDecorator, error) {
	node, err := htmlquery.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	return NewNodeFromHTML(node)
}

func NewNodeFromHTML(node *html.Node) (*NodeDecorator, error) {
	if node == nil {
		return nil, errors.New("nil html node")
	}
	return &NodeDecorator{*node}, nil
}

func (n *NodeDecorator) ToNode() *html.Node {
	return &n.Node
}

func (n *NodeDecorator) GetString() string {
	return htmlquery.InnerText(n.ToNode())
}

func (n *NodeDecorator) GetInnerNodeList(expr string) ([]*NodeDecorator, error) {
	innerNodes, err := htmlquery.QueryAll(n.ToNode(), expr)
	if err != nil {
		return nil, err
	}
	var nodeList []*NodeDecorator
	for _, node := range innerNodes {
		if n, err := NewNodeFromHTML(node); err != nil {
			return nodeList, err
		} else {
			nodeList = append(nodeList, n)
		}
	}
	return nodeList, nil
}

func (n *NodeDecorator) GetInnerNode(expr string) (*NodeDecorator, error) {
	node, err := htmlquery.Query(n.ToNode(), expr)
	if err != nil {
		return nil, err
	}
	return NewNodeFromHTML(node)
}

func (n *NodeDecorator) GetInnerString(expr string) (string, error) {
	node, err := htmlquery.Query(n.ToNode(), expr)
	if err != nil {
		return "", err
	}
	wrapped, err := NewNodeFromHTML(node)
	if err != nil {
		return "", err
	}
	return wrapped.GetString(), nil
}
