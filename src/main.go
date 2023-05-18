package main

import (
	"image/color"
	"sort"
	"time"

	"os"

	"strconv"

	"fmt"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
	"github.com/google/uuid"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide a port number")
		return
	}

	algorithm := args[0]
	if algorithm != "kleppmann" && algorithm != "lumina" {
		fmt.Println("Please provide a valid algorithm name: (kleppmann, lumina)")
		return
	}

	numPeers, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Println("Please provide a valid number of peers")
		return
	}
	port, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		fmt.Println("Please provide a valid port number")
		return
	}

	var tree treeinterface.Tree[string]
	if algorithm == "kleppmann" {
		tree = treeinterface.NewKTree[string](connection.NewTCPProvider(int(numPeers), int(port)))
	} else {
		tree = treeinterface.NewLTree[string](connection.NewTCPProvider(int(numPeers), int(port)), false)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("My App")

	var treeDisplay *widget.Tree
	treeDisplay = &widget.Tree{
		Root: tree.Root().String(),
		ChildUIDs: func(uid string) []string {
			id, err := uuid.Parse(uid)
			if err != nil {
				return nil
			}
			chld, err := tree.GetChildren(id)
			if err != nil {
				return nil
			}
			var res []string
			for _, v := range chld {
				res = append(res, v.String())
			}
			sort.Slice(res, func(i, j int) bool {
				return strings.Compare(res[i], res[j]) < 0
			})
			return res
		},
		IsBranch: func(uid string) bool {
			id, err := uuid.Parse(uid)
			if err != nil {
				return false
			}
			chld, err := tree.GetChildren(id)
			if err != nil {
				return false
			}
			return len(chld) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			hbox := container.NewHBox()

			id := widget.NewLabel("ID")
			hbox.Add(id)
			hbox.Add(widget.NewEntry())
			hbox.Add(widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
				println("clicked")
			}))
			hbox.Add(widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
				println("clicked")
			}))
			delete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				println("clicked")
			})
			// make delete button red
			delete.Importance = widget.DangerImportance
			hbox.Add(delete)

			return hbox
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			// if uid == tree.Root().String() {
			// 	hbox.Objects[0].(*widget.Label).SetText("Root")
			// 	hbox.Objects[0].(*widget.Label).TextStyle.Bold = true
			// 	hbox.Objects[2].(*widget.Button).OnTapped = func() {
			// 		val := hbox.Objects[1].(*widget.Entry).Text
			// 		//val, _ := str.Get()
			// 		tree.Insert(tree.Root(), val)
			// 		hbox.Objects[1].(*widget.Entry).SetText("")
			// 		treeDisplay.Refresh()
			// 	}
			// 	hbox.Objects[4].(*widget.Button).Hide()
			// 	return
			// }

			id, err := uuid.Parse(uid)
			if err != nil {
				return
			}
			node, err := tree.Get(id)
			if err != nil {
				return
			}
			meta := node.Metadata()
			hbox.Objects[0].(*widget.Label).SetText("'" + meta + "'")

			//str := binding.NewString()
			//hbox.Objects[1].(*widget.Entry).Bind(str)
			hbox.Objects[2].(*widget.Button).OnTapped = func() {
				val := hbox.Objects[1].(*widget.Entry).Text
				if val == "" {
					return
				}
				//val, _ := str.Get()
				tree.Insert(id, val)
				hbox.Objects[1].(*widget.Entry).SetText("")
				treeDisplay.Refresh()
			}
			hbox.Objects[3].(*widget.Button).OnTapped = func() {
				myWindow.Clipboard().SetContent(uid)
			}

			hbox.Objects[4].(*widget.Button).OnTapped = func() {
				tree.Delete(id)
				treeDisplay.Refresh()
			}
		},
		OnSelected: func(id string) {
			println("selected", id)
			myWindow.Clipboard().SetContent(id)
		},
	}

	//refresh tree every second to update the view
	left := container.NewVBox()
	menuLabel := widget.NewLabel("Menu")

	left.Add(menuLabel)

	title := canvas.NewText("Tree CRDTs - "+strings.Title(algorithm)+" ("+strconv.FormatInt(port, 10)+")", color.White)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true
	title.TextSize = 24

	// Bottom will have a port input and a connect button
	//portLabel := widget.NewLabel("Port: ")
	portInput := widget.NewEntry()
	portInput.SetPlaceHolder("Port number")
	portInput.Resize(fyne.NewSize(200, 50))

	connectButton := widget.NewButton("Connect", func() {
		port, err := strconv.ParseInt(portInput.Text, 10, 64)
		if err != nil {
			fmt.Println("Please provide a valid port number")
			return
		}
		if port < 0 || port > 65535 {
			fmt.Println("Please provide a valid port number")
			return
		}
		address := "localhost:" + portInput.Text
		tree.ConnectionProvider().Connect(address)
		treeDisplay.Refresh()
	})
	// make it green
	connectButton.Importance = widget.HighImportance

	disconnectButton := widget.NewButton("Disconnect from all", func() {
		tree.ConnectionProvider().CloseAll()
	})
	// make it red
	disconnectButton.Importance = widget.WarningImportance

	moveFromEntry := widget.NewEntry()
	moveFromEntry.SetPlaceHolder("Move from")
	moveToEntry := widget.NewEntry()
	moveToEntry.SetPlaceHolder("Move to below")
	doMove := widget.NewButton("Move", func() {
		from, err := uuid.Parse(moveFromEntry.Text)
		if err != nil {
			fmt.Println("Please provide a valid ID")
			return
		}
		to, err := uuid.Parse(moveToEntry.Text)
		if err != nil {
			fmt.Println("Please provide a valid ID")
			return
		}
		tree.Move(from, to)
		// clear the entries
		moveFromEntry.SetText("")
		moveToEntry.SetText("")
		treeDisplay.Root = ""
		treeDisplay.Refresh()
		treeDisplay.Root = tree.Root().String()
		treeDisplay.Refresh()
	})

	move := container.NewGridWithColumns(3, moveFromEntry, moveToEntry, doMove)
	connection := container.NewGridWithColumns(3, portInput, connectButton, disconnectButton)

	bottom := container.NewVBox(move, connection)

	right := container.NewVBox()

	content := container.NewBorder(title, bottom, nil, right, treeDisplay)

	go func() {
		for range time.Tick(time.Second) {
			treeDisplay.Refresh()
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			connections := tree.ConnectionProvider().GetPeerAddrs()
			right.Objects = nil
			cons := widget.NewLabel("Connections")
			cons.TextStyle.Bold = true
			right.Add(cons)
			for _, v := range connections {
				address := v.String()
				port := strings.Split(address, ":")[1]
				right.Add(canvas.NewText(port, color.White))
			}
		}
	}()
	myWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyF5 {
			treeDisplay.Root = ""
			treeDisplay.Refresh()
			treeDisplay.Root = tree.Root().String()
			treeDisplay.Refresh()
		}
	})
	myWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyF4 {
			treeDisplay.OpenAllBranches()
		}
	})
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 400))
	myWindow.ShowAndRun()
}
