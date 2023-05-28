package main

import (
	"github.com/rivo/tview"
	"github.com/strider2038/algos/search_in_trees/bitmap_search"
)

type TreeViewer struct {
	classifiers *bitmap_search.Tree
	height      int
	maxHeight   int
}

func renderTree(classifiers *bitmap_search.Tree, parent *tview.TreeNode, search string, maxHeight int) {
	viewer := &TreeViewer{classifiers: classifiers, maxHeight: maxHeight}
	viewer.filterNodes(parent, "", search, 0)
}

func (v *TreeViewer) filterNodes(parent *tview.TreeNode, code, search string, level int) {
	var nodes []*bitmap_search.Node
	if code == "" {
		nodes = v.classifiers.Filter(search)
	} else {
		nodes = v.classifiers.FilterAt(code, search)
	}

	parent.ClearChildren()
	for _, node := range nodes {
		child := tview.NewTreeNode(node.Value.Code + " - " + node.Value.Title)
		parent.AddChild(child)
		v.height++
		v.filterNodes(child, node.Value.Code, search, level+1)

		if v.height > v.maxHeight {
			return
		}
	}
}
