package module

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"strings"
)

// NodeDecorator may become nil, all func must prevent nil call/*
type NodeDecorator struct {
	html.Node
}

func NewNodeFromBytes(data []byte) (*NodeDecorator, error) {
	node, err := htmlquery.Parse(strings.NewReader(string(data)))
	return NewNodeFromHTML(node), err
}

func NewNodeFromHTML(node *html.Node) *NodeDecorator {
	if node == nil {
		return nil
	}
	return &NodeDecorator{*node}
}

func (n *NodeDecorator) ToNode() *html.Node {
	if n == nil {
		return nil
	}
	return &n.Node
}

func (n *NodeDecorator) GetString() string {
	if n == nil {
		return ""
	}
	return htmlquery.InnerText(n.ToNode())
}

func (n *NodeDecorator) GetInnerNodeList(expr string) []*NodeDecorator {
	if n == nil {
		return nil
	}
	innerNodes, err := htmlquery.QueryAll(n.ToNode(), expr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var nodeList []*NodeDecorator
	for _, node := range innerNodes {
		nodeList = append(nodeList, NewNodeFromHTML(node))
	}
	return nodeList
}

func (n *NodeDecorator) GetInnerNode(expr string) *NodeDecorator {
	if n == nil {
		return nil
	}
	node, err := htmlquery.Query(n.ToNode(), expr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return NewNodeFromHTML(node)
}
