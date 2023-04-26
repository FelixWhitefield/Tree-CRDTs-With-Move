package treeinterface

import (
	"testing"
	"time"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"github.com/google/uuid"
)

func TestKTreeOperationTransmits(t *testing.T) {
	tree1 := NewKTree[string](connection.NewTCPProvider(1, 1111))
	tree2 := NewKTree[string](connection.NewTCPProvider(1, 1112))

	tree1.ConnectionProvider().Connect("localhost:1112")

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
}

func TestKTreeCycleMove(t *testing.T) {
	tree1 := NewKTree[string](connection.NewTCPProvider(1, 1211))
	tree2 := NewKTree[string](connection.NewTCPProvider(1, 1212))

	tree1.ConnectionProvider().Connect("localhost:1212")

	id1, _ := tree1.Insert(tree1.Root(), "meta") // Add two nodes
	id2, _ := tree1.Insert(tree1.Root(), "meta1")

	time.Sleep(1 * time.Second) // Time for communication to occur
	

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

func contains(arr []uuid.UUID, val uuid.UUID) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}