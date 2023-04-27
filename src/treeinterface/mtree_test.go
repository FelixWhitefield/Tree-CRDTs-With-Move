package treeinterface

import (
	"testing"
	"time"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
)

func TestMTreeOperationTransmits(t *testing.T) {
	tree1 := NewMTree[string](connection.NewTCPProvider(2, 2221))
	tree2 := NewMTree[string](connection.NewTCPProvider(2, 2222))

	tree1.ConnectionProvider().Connect("localhost:2222")

	nodeId, err := tree1.Insert(tree1.Root(), "meta")
	if err != nil {
		t.Errorf("Could not insert node")
	}

	time.Sleep(1 * time.Second) // Time for communication to occur

	if id, _ := tree2.Get(nodeId); id == nil {
		t.Errorf("Node was not inserted")
	}

	nodeId1, err := tree1.Insert(tree1.Root(), "meta1")
	if err != nil {
		t.Errorf("Could not insert node1")
	}

	nodeId2, err := tree1.Insert(tree1.Root(), "meta2")
	if err != nil {
		t.Errorf("Could not insert node2")
	}

	nodeId3, err := tree1.Insert(tree1.Root(), "meta3")
	if err != nil {
		t.Errorf("Could not insert node3")
	}

	time.Sleep(1 * time.Second) // Time for communication to occur

	rootChildren, err := tree2.GetChildren(tree2.Root())
	if err != nil {
		t.Errorf("Could not get children")
	}

	if len(rootChildren) != 4 && !contains(rootChildren, nodeId1) && !contains(rootChildren, nodeId2) && !contains(rootChildren, nodeId3) && !contains(rootChildren, nodeId) {
		t.Errorf("Expected 4 children, got %d", len(rootChildren))
	}

	tree3 := NewMTree[string](connection.NewTCPProvider(2, 2223))
	tree3.ConnectionProvider().Connect("localhost:2221")
	tree3.ConnectionProvider().Connect("localhost:2222")
	time.Sleep(2 * time.Second) // Time for communication to occur

	rootChildren3, err := tree3.GetChildren(tree3.Root())
	if err != nil {
		t.Errorf("Could not get children")
	}

	if len(rootChildren3) != 4 && !contains(rootChildren3, nodeId1) && !contains(rootChildren3, nodeId2) && !contains(rootChildren3, nodeId3) && !contains(rootChildren3, nodeId) {
		t.Errorf("Expected 4 children, got %d", len(rootChildren3))
	}

}

func TestMTreeCycleMove(t *testing.T) {
	tree1 := NewMTree[string](connection.NewTCPProvider(1, 3223))
	tree2 := NewMTree[string](connection.NewTCPProvider(1, 3224))

	tree1.ConnectionProvider().Connect("localhost:3224")

	id1, err := tree1.Insert(tree1.Root(), "meta") // Add two nodes
	if err != nil {
		t.Errorf("Could not insert node")
	}
	id2, err := tree1.Insert(tree1.Root(), "meta1")
	if err != nil {
		t.Errorf("Could not insert node")
	}

	time.Sleep(1 * time.Second) // Time for communication to occur

	//check both trees have the nodes
	if _, err := tree2.Get(id1); err != nil {
		t.Errorf("Node was not inserted %v, with id: %v", err, id1)
	}
	if _, err := tree2.Get(id2); err != nil {
		t.Errorf("Node was not inserted %v", err)
	}
	children2, _ := tree2.GetChildren(tree2.Root())
	if len(children2) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children2))
	}

	// Tree is now:
	//     Root
	//    /    \
	//  id1    id2

	// Move nodes in a cycle
	tree1.Move(id1, id2)
	tree2.Move(id2, id1)

	// Tree will either be:
	//     Root
	//    /
	//  id2
	//  id1
	// or
	//     Root
	//    /
	//  id1
	//  id2

	time.Sleep(1 * time.Second) // Time for communication to occur

	// Check that states are the same after move
	if !tree1.crdt.Equals(tree2.crdt) {
		t.Errorf("States are not the same after move")
	}
	children, _ := tree1.GetChildren(tree1.Root())
	layer2, _ := tree1.GetChildren(children[0])

	if len(children) != 1 && len(layer2) != 1 {
		t.Errorf("Expected 1 child, got %d and %d", len(children), len(layer2))
	}
}
