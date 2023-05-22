package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/strider2038/algos/search_in_trees/bitmap_search"
)

func main() {
	app := tview.NewApplication()

	classifiers, err := bitmap_search.LoadFromFile("./testdata/classifiers/okved.csv")
	if err != nil {
		log.Fatal(err)
	}

	root := tview.NewTreeNode("ОКВЭД")
	for _, node := range classifiers.Nodes {
		root.AddChild(tview.NewTreeNode(node.Value.Code + " - " + node.Value.Title))
	}
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	inputField := tview.NewInputField().
		SetLabel("Подбор классификатора: ").
		SetFieldWidth(100).
		SetChangedFunc(func(search string) {
			nodes := classifiers.Filter(search)

			root.ClearChildren()
			for _, node := range nodes {
				root.AddChild(tview.NewTreeNode(node.Value.Code + " - " + node.Value.Title))
			}
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
