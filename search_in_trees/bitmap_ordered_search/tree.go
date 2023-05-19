package bitmap_ordered_search

import (
	"errors"

	"github.com/strider2038/algos/prefix_trees/byte_suffix_trie"
)

var errKeywordNotFound = errors.New("keyword not found")

type Classifier struct {
	Code  string
	Title string
}

type Tree struct {
	Nodes []*Node

	keywords byte_suffix_trie.Array[int]
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
	keywords := parseKeywords(search, true)
	if len(keywords) == 0 {
		keywords = parseKeywords(search, false)
	}

	// для совпадения условий поиска необходимо, чтобы все группы
	// ключевых слов совпали
	ranges := make([]indexRange, 0, len(keywords))
	for _, keyword := range keywords {
		// в качестве совпадения может подойти любое
		// из ключевых слов, найденных по префиксу
		r := indexRange{}
		k := []byte(keyword)
		if _, index, ok := t.keywords.FindFirstByPrefix(k); ok {
			r.First = index
		} else {
			// попалось слово, которого нет в словаре
			return nil
		}
		if _, index, ok := t.keywords.FindLastByPrefix(k); ok {
			r.Last = index
		} else {
			// попалось слово, которого нет в словаре
			return nil
		}

		ranges = append(ranges, r)
	}

	nodes := make([]*Node, 0)
	for _, node := range t.Nodes {
		if node.matches(ranges) {
			nodes = append(nodes, &Node{Value: node.Value})
		}
	}

	return nodes
}

func (t *Tree) Bytes() int {
	s := 0
	for _, node := range t.Nodes {
		s += node.bytes()
	}
	treeBytes, treeNodes := t.keywords.Size()

	return s + treeBytes + treeNodes*8
}

func (t *Tree) newNode(code string, title string) (*Node, error) {
	node := &Node{
		Value:   Classifier{Code: code, Title: title},
		indices: map[int]uint64{},
	}

	keywords := parseKeywords(title, true)
	for _, keyword := range keywords {
		if i, exist := t.keywords.Find([]byte(keyword)); exist {
			node.indices.Set(i)
		} else {
			return nil, errKeywordNotFound
		}
	}

	return node, nil
}

func (t *Tree) reindex() {
	for _, node := range t.Nodes {
		node.reindex()
	}
}

type Node struct {
	Value    Classifier
	Children []*Node

	indices bitmap
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

func (n *Node) reindex() {
	for _, child := range n.Children {
		child.reindex()
		n.indices = n.indices.Or(child.indices)
	}
}

func (n *Node) matches(ranges []indexRange) bool {
	for _, r := range ranges {
		if !n.containsOneOf(r) {
			return false
		}
	}

	return true
}

func (n *Node) containsOneOf(indices indexRange) bool {
	for index := indices.First; index <= indices.Last; index++ {
		if n.indices.IsSet(index) {
			return true
		}
	}

	return false
}

func (n *Node) bytes() int {
	s := 0
	for _, child := range n.Children {
		s += child.bytes()
	}

	return s + len(n.Value.Code) + len(n.Value.Title) + len(n.indices)*8
}

type indexRange struct {
	First int
	Last  int
}
