package treeinterface

import (
	"log"
	"testing"
	"time"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
)

func TestMTreeOperationTransmits(t *testing.T) {
	tree1 := NewMTree[string](connection.NewTCPProvider(1, 2221))
	tree2 := NewMTree[string](connection.NewTCPProvider(1, 2222))

	tree1.ConnectionProvider().Connect("localhost:2222")

	nodeId, err := tree1.Insert(tree1.Root(), "meta")
	if err != nil {
		t.Errorf("Could not insert node")
	}

	time.Sleep(1 * time.Second) // Time for communication to occur

	if id, _ := tree2.Get(nodeId); id != nil {
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
		log.Println(rootChildren)
		tree1Children, _ := tree1.GetChildren(tree1.Root())
		log.Println(tree1Children)
		t.Errorf("Expected 4 children, got %d", len(rootChildren))
	}
}
