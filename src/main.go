package main 

import (
	ti "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
	"fmt"
	"time"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
)

func main() {
		//var err error
	
		// mtree := ti.NewLTree[string](connection.NewTCPProvider(2, 1122))
		// m2tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1123))
		// m3tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1124))
	
		// m2tree.ConnectionProvider().Connect("localhost:1122")
		// m3tree.ConnectionProvider().Connect("localhost:1122")
	
		mtree := ti.NewLTree[string](connection.NewTCPProvider(2, 1122), true)
		m2tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1123), true)
		m3tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1124), true)
	
		
		m2tree.ConnectionProvider().Connect("localhost:1122")
		m3tree.ConnectionProvider().Connect("localhost:1122")
		
		time.Sleep(1 * time.Second)
	
		id1, _ := mtree.Insert(mtree.Root(), "root")
		id2, _ := mtree.Insert(mtree.Root(), "root2")
		id3, _ := mtree.Insert(mtree.Root(), "root3")

		fmt.Printf("ID1: %v, ID2: %v, ID3: %v\n", id1, id2, id3)

		time.Sleep(1 * time.Second)
		
		mtree.Move(id1, id2)
		m2tree.Delete(id2)
		m3tree.Move(id2, id1)
		mtree.Move(id3, id2)
		m2tree.Move(id3, id2)

		time.Sleep(1 * time.Second)
	
		nodes11, _ := mtree.GetChildren(mtree.Root())
		nodes22, _ := m2tree.GetChildren(m2tree.Root())
		nodes33, _ := m3tree.GetChildren(m3tree.Root())

		for _, node := range nodes11 {
			nodeChildren, _ := mtree.GetChildren(node)
			fmt.Printf("1 Node: %v  \n", nodeChildren)
		}

		for _, node := range nodes22 {
			nodeChildren, _ := m2tree.GetChildren(node)
			fmt.Printf("2 Node: %v  \n", nodeChildren)
		}

		for _, node := range nodes33 {
			nodeChildren, _ := m3tree.GetChildren(node)
			fmt.Printf("3 Node: %v  \n", nodeChildren)
		}

		fmt.Println(mtree)
	
		fmt.Printf("Nodes: %v  \n", nodes11)
		fmt.Printf("Nodes: %v  \n", nodes22)
		fmt.Printf("Nodes: %v  \n", nodes33)
		
		
		//fmt.Println(m2tree.Equals(mtree))
}