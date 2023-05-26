package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/strider2038/algos/search_in_trees/bitmap_search"
)

var maxHeight = 30

func main() {
	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		_, h := screen.Size()
		if h <= 10 {
			h = 10
		}
		maxHeight = h

		return false
	})

	classifiers, err := bitmap_search.LoadFromFile("./var/ksr.csv")
	//classifiers, err := bitmap_search.LoadFromFile("./testdata/classifiers/okved.csv")
	if err != nil {
		log.Fatal(err)
	}

	root := tview.NewTreeNode("ОКВЭД")
	renderTree(classifiers, root, "")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	inputField := tview.NewInputField().
		SetLabel("Подбор классификатора: ").
		SetFieldWidth(100).
		SetChangedFunc(func(text string) {
			renderTree(classifiers, root, text)
		}).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(tree, 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		log.Fatal(err)
	}
}

type TreeViewer struct {
	classifiers *bitmap_search.Tree
	height      int
}

func renderTree(classifiers *bitmap_search.Tree, parent *tview.TreeNode, search string) {
	viewer := &TreeViewer{classifiers: classifiers}
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
	for i, node := range nodes {
		child := tview.NewTreeNode(node.Value.Code + " - " + node.Value.Title)
		parent.AddChild(child)
		v.height++
		v.filterNodes(child, node.Value.Code, search, level+1)

		if v.height > maxHeight {
			return
		}

		switch level {
		case 1:
			if i >= 5 {
				if len(nodes) > 5 {
					parent.AddChild(tview.NewTreeNode("..."))
				}
				return
			}
		case 2:
			if i >= 3 {
				if len(nodes) > 3 {
					parent.AddChild(tview.NewTreeNode("..."))
				}
				return
			}
		default:
			if i > 1 {
				if len(nodes) > 1 {
					parent.AddChild(tview.NewTreeNode("..."))
				}
				return
			}
		}
	}
}
