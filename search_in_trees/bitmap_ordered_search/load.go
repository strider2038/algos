package bitmap_ordered_search

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/strider2038/algos/prefix_trees/byte_suffix_trie"
)

func LoadFromFile(filename string) (*Tree, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("read from file: %w", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read from CSV: %w", err)
	}

	keywords := createOrderedKeywords(rows)
	loader := newTreeLoader(rows)
	loader.tree.keywords = keywords
	if err := loader.fill(); err != nil {
		return nil, err
	}
	loader.tree.reindex()

	return &loader.tree, nil
}

type treeLoader struct {
	rows [][]string
	tree Tree
}

func newTreeLoader(rows [][]string) *treeLoader {
	return &treeLoader{
		rows: rows,
		tree: newTree(len(rows)),
	}
}

func (l *treeLoader) fill() error {
	for i, row := range l.rows {
		node, err := l.tree.newNode(row[1], row[2])
		if err != nil {
			return fmt.Errorf("error at %d (%s): %w", i, row[1], err)
		}
		if row[0] == "" {
			l.tree.Nodes = append(l.tree.Nodes, node)
		} else {
			parent := l.tree.nodesByCodes[row[0]]
			parent.Children = append(parent.Children, node)
		}
		l.tree.nodesByCodes[node.Value.Code] = node
	}

	return nil
}

func createOrderedKeywords(rows [][]string) byte_suffix_trie.Array[int] {
	dictionary := byte_suffix_trie.Array[struct{}]{}
	for _, row := range rows {
		keywords := parseKeywords(row[2], true)
		for _, keyword := range keywords {
			dictionary.Put([]byte(keyword), struct{}{})
		}
	}
	index := 1
	keywords := byte_suffix_trie.Array[int]{}
	_ = dictionary.Walk(func(key []byte, value struct{}) error {
		keywords.Put(key, index)
		index++

		return nil
	})

	return keywords
}
