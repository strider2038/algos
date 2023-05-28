package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/strider2038/algos/search_in_trees/bitmap_search"
)

const (
	minHeight     = 10
	defaultHeight = 30

	selectClassifierPageName = "classifier_select"
	filterClassifierPageName = "classifier_filter"
)

var sourceFiles = map[string]string{
	"КСР":   "./../../var/ksr.csv",
	"ОКВЭД": "./../../testdata/classifiers/okved.csv",
}

type UI struct {
	app      *tview.Application
	pages    *tview.Pages
	treeView *tview.TreeView

	classifiers *bitmap_search.Tree

	maxHeight int
}

func NewUI() *UI {
	ui := &UI{maxHeight: defaultHeight}
	ui.app = tview.NewApplication()
	ui.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		_, h := screen.Size()
		if h <= minHeight {
			h = minHeight
		}
		ui.maxHeight = h

		return false
	})
	ui.pages = tview.NewPages()

	// first page
	classifierSelector := tview.NewList()
	classifierSelector.AddItem("ОКВЭД", "", '1', func() {
		ui.loadTree("ОКВЭД")
		ui.pages.SwitchToPage(filterClassifierPageName)
	})
	classifierSelector.AddItem("КСР", "", '2', func() {
		ui.loadTree("КСР")
		ui.pages.SwitchToPage(filterClassifierPageName)
	})
	classifierSelector.AddItem("Выйти", "", 'q', func() {
		ui.app.Stop()
	})

	ui.pages.AddPage(selectClassifierPageName, classifierSelector, false, true)

	// second page
	ui.treeView = tview.NewTreeView()

	classifierFilterInput := tview.NewInputField().
		SetLabel("Подбор классификатора: ").
		SetFieldWidth(100).
		SetChangedFunc(func(text string) {
			ui.renderTree(text)
		}).
		SetDoneFunc(func(key tcell.Key) {
			ui.app.Stop()
		})

	classifierFilterScreen := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(classifierFilterInput, 1, 1, true).
		AddItem(ui.treeView, 0, 1, false)

	ui.pages.AddPage(filterClassifierPageName, classifierFilterScreen, true, false)

	return ui
}

func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.pages, true).SetFocus(ui.pages).Run()
}

func (ui *UI) loadTree(classifierName string) {
	classifiers, err := bitmap_search.LoadFromFile(sourceFiles[classifierName])
	if err != nil {
		log.Fatal(err)
	}
	ui.classifiers = classifiers

	root := tview.NewTreeNode(classifierName)
	ui.treeView.SetRoot(root).SetCurrentNode(root)
	ui.renderTree("")
}

func (ui *UI) renderTree(search string) {
	renderTree(ui.classifiers, ui.treeView.GetRoot(), search, ui.maxHeight)
}
