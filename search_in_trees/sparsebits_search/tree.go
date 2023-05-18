package sparsebits_search

import (
	"github.com/strider2038/algos/bitmaps"
	"github.com/strider2038/algos/prefix_trees/byte_suffix_trie"
)

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
	and := make([][]int, 0, len(keywords))
	for _, keyword := range keywords {
		// в качестве совпадения может подойти любое
		// из ключевых слов, найденных по префиксу
		or := make([]int, 0, 1)
		_ = t.keywords.WalkPrefix([]byte(keyword), func(key []byte, index int) error {
			or = append(or, index)
			return nil
		})

		// попалось слово, которого нет в словаре
		if len(or) == 0 {
			return nil
		}
		and = append(and, or)
	}

	nodes := make([]*Node, 0)
	for _, node := range t.Nodes {
		if node.matches(and) {
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

func (t *Tree) newNode(code string, title string) *Node {
	node := &Node{
		Value: Classifier{Code: code, Title: title},
	}

	keywords := parseKeywords(title, true)
	index := t.keywords.Count() + 1
	for _, keyword := range keywords {
		if i, exist := t.keywords.Find([]byte(keyword)); exist {
			node.indices.Set(i)
		} else {
			t.keywords.Put([]byte(keyword), index)
			node.indices.Set(index)
			index++
		}
	}

	return node
}

func (t *Tree) reindex() {
	for _, node := range t.Nodes {
		node.reindex()
	}
}

type Node struct {
	Value    Classifier
	Children []*Node

	indices bitmaps.SparseBitmap64p3
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

func (n *Node) matches(and [][]int) bool {
	for _, or := range and {
		if !n.containsOneOf(or) {
			return false
		}
	}

	return true
}

func (n *Node) containsOneOf(indices []int) bool {
	for _, index := range indices {
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

	return s + len(n.Value.Code) + len(n.Value.Title) + n.indices.Bytes()
}
