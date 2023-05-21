package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	root := tview.NewTreeNode("ОКВЭД")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	inputField := tview.NewInputField().
		SetLabel("Подбор классификатора: ").
		SetFieldWidth(100).
		SetChangedFunc(func(text string) {
			root.ClearChildren()
			root.AddChild(tview.NewTreeNode(text))
			root.AddChild(tview.NewTreeNode(text + " 1"))
			root.AddChild(tview.NewTreeNode(text + " 2"))
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
