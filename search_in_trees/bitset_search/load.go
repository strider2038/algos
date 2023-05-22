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

	loader := newTreeLoader(rows)
	loader.fill()
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

func (l *treeLoader) fill() {
	for _, row := range l.rows {
		node := l.tree.newNode(row[1], row[2])
		if row[0] == "" {
			l.tree.Nodes = append(l.tree.Nodes, node)
		} else {
			parent := l.tree.nodesByCodes[row[0]]
			parent.Children = append(parent.Children, node)
		}
		l.tree.nodesByCodes[node.Value.Code] = node
	}
}
