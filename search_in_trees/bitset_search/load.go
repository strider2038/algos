package bitset_search

import (
	"encoding/csv"
	"fmt"
	"os"
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

	loader := treeLoader{rows: rows}
	loader.fill()
	loader.tree.reindex()

	return &loader.tree, nil
}

type treeLoader struct {
	rows   [][]string
	offset int
	tree   Tree
}

func (l *treeLoader) fill() {
	for l.offset < len(l.rows) {
		row := l.rows[l.offset]
		if row[1] == "" {
			node := l.tree.newNode(row[0], row[2])
			l.tree.Nodes = append(l.tree.Nodes, node)
			l.offset++
			if l.offset >= len(l.rows) {
				break
			}
			l.fillChildren(node)
		}
	}
}

func (l *treeLoader) fillChildren(parent *Node) {
	code := l.rows[l.offset][1]

	for l.offset < len(l.rows) {
		row := l.rows[l.offset]
		if row[1] == "" {
			return
		}

		node := l.tree.newNode(row[1], row[2])
		parent.Children = append(parent.Children, node)
		l.offset++
		if l.offset >= len(l.rows) {
			break
		}
		if len(l.rows[l.offset][1]) > len(code) {
			l.fillChildren(node)
		}
		if len(l.rows[l.offset][1]) < len(code) {
			break
		}
	}
}
