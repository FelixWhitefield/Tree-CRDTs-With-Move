package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
	"github.com/google/uuid"
)

type TreeNode struct {
	id   string // unique id
	name string // name to display
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("My App")

	ktree := treeinterface.NewKTree[string](connection.NewTCPProvider(1, 1111))

	tree := &widget.Tree{
		Root: ktree.Root().String(),
		ChildUIDs: func(uid string) []string {
			id, err := uuid.Parse(uid)
			if err != nil {
				return nil
			}
			chld, err := ktree.GetChildren(id)
			if err != nil {
				return nil
			}
			var res []string
			for _, v := range chld {
				res = append(res, v.String())
			}
			return res
		},
		IsBranch: func(uid string) bool {
			return true
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			hbox := container.NewHBox()
			hbox.Add(widget.NewLabel("Hello"))
			hbox.Add(widget.NewButton("Click me", func() {
				println("clicked")
			}))
			return hbox
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			hbox.Objects = nil

			id, err := uuid.Parse(uid)
			if err != nil {
				return
			}
			node, err := ktree.Get(id)
			if err != nil {
				return
			}
			meta := node.Metadata()
			hbox.Add(widget.NewLabel(meta))
			delete := widget.NewButton("Delete", func() {
				err := ktree.Delete(id)
				fmt.Println(ktree.GetChildren(ktree.Root()))
				fmt.Println(err)
				fmt.Println("deleted", id.String())
			})
			// make delete button red
			delete.Importance = widget.DangerImportance
			hbox.Add(delete)
			hbox.Add(widget.NewButtonWithIcon("copy", theme.ContentCopyIcon(), func() {
				myWindow.Clipboard().SetContent(uid)
			}))
		},
		OnSelected: func(id string) {
			println("selected", id)
			myWindow.Clipboard().SetContent(id)
		},
	}

	ktree.Insert(ktree.Root(), "Hello")
	tree.Refresh()

	//refresh tree every second to update the view
	go func() {
		for range time.Tick(time.Second) {
			tree.Refresh()
		}
	}()

	content := container.New(layout.NewGridLayout(1), tree)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 400))
	myWindow.ShowAndRun()
}
