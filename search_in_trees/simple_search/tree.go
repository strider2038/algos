package simple_search

import (
	"strings"

	"github.com/kljensen/snowball/russian"
)

type Classifier struct {
	Code  string
	Title string
}

type Tree struct {
	Nodes []*Node
}

func (t *Tree) FindByCode(code string) (*Node, bool) {
	for _, node := range t.Nodes {
		if n, ok := node.findByCode(code); ok {
			return n, ok
		}
	}

	return nil, false
}

func (t *Tree) Filter(search string) []*Node {
	keywords := strings.Split(strings.ToLower(search), " ")
	for i := range keywords {
		keywords[i] = russian.Stem(keywords[i], true)
	}

	nodes := make([]*Node, 0)
	for _, node := range t.Nodes {
		if node.contains(keywords) {
			nodes = append(nodes, &Node{Value: node.Value})
		}
	}

	return nodes
}

type Node struct {
	Value    Classifier
	Children []*Node
}

func (n *Node) findByCode(code string) (*Node, bool) {
	if n.Value.Code == code {
		return n, true
	}

	for _, child := range n.Children {
		if c, ok := child.findByCode(code); ok {
			return c, true
		}
	}

	return nil, false
}

func (n *Node) contains(keywords []string) bool {
	if contains(n.Value.Title, keywords) {
		return true
	}

	for _, child := range n.Children {
		if child.contains(keywords) {
			return true
		}
	}

	return false
}

func contains(s string, keywords []string) bool {
	s = strings.ToLower(s)

	for _, keyword := range keywords {
		if !strings.Contains(s, keyword) {
			return false
		}
	}

	return true
}
