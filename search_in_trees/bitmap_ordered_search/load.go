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
	loader := treeLoader{rows: rows}
	loader.tree.keywords = keywords
	if err := loader.fill(); err != nil {
		return nil, err
	}
	loader.tree.reindex()

	return &loader.tree, nil
}

type treeLoader struct {
	rows   [][]string
	offset int
	tree   Tree
}

func (l *treeLoader) fill() error {
	for l.offset < len(l.rows) {
		row := l.rows[l.offset]
		if row[1] == "" {
			node, err := l.tree.newNode(row[0], row[2])
			if err != nil {
				return fmt.Errorf("error at %d (%s): %w", l.offset, row[0], err)
			}
			l.tree.Nodes = append(l.tree.Nodes, node)
			l.offset++
			if l.offset >= len(l.rows) {
				break
			}
			if err := l.fillChildren(node); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *treeLoader) fillChildren(parent *Node) error {
	code := l.rows[l.offset][1]

	for l.offset < len(l.rows) {
		row := l.rows[l.offset]
		if row[1] == "" {
			return nil
		}

		node, err := l.tree.newNode(row[1], row[2])
		if err != nil {
			return fmt.Errorf("error at %d (%s): %w", l.offset, row[1], err)
		}
		parent.Children = append(parent.Children, node)
		l.offset++
		if l.offset >= len(l.rows) {
			break
		}
		if len(l.rows[l.offset][1]) > len(code) {
			if err := l.fillChildren(node); err != nil {
				return err
			}
		}
		if len(l.rows[l.offset][1]) < len(code) {
			break
		}
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
