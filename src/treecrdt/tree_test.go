package treecrdt

import (
	"github.com/google/uuid"
	"testing"
)

func TestTreeAddNodes(t *testing.T) {
	tree := NewTree[string]()

	childId1 := uuid.New()
	tree.Add(childId1, NewTreeNode(RootUUID, "test1")) // Test adding a node

	children, _  := tree.GetChildren(RootUUID)
	if len(children) == 1 && children[0] != childId1 {
		t.Errorf("Tree.GetChildren() returned %v, expected %v", children, []uuid.UUID{childId1})
	}

	childId2 := uuid.New()
	tree.Add(childId2, NewTreeNode(RootUUID, "test2")) // Test adding another node

	children, _  = tree.GetChildren(RootUUID)
	if len(children) == 2 && !contains(children, childId1) && !contains(children, childId2) {
		t.Errorf("Tree.GetChildren() returned %v, expected %v", children, []uuid.UUID{childId1, childId2})
	}

	childId3 := uuid.New()
	tree.Add(childId3, NewTreeNode(childId1, "test3")) // Test adding a node under another node

	children, _  = tree.GetChildren(childId1)
	if len(children) == 1 && children[0] != childId3 {
		t.Errorf("Tree.GetChildren() returned %v, expected %v", children, []uuid.UUID{childId3})
	}
	
	if tree.GetNode(childId3).PrntID != childId1 {
		t.Errorf("Tree.GetNode().PrntID returned %v, expected %v", tree.GetNode(childId3).PrntID, childId1)
	}
}

func TestTreeRemoveNodes(t *testing.T) {
	tree := NewTree[string]()

	childId1 := uuid.New()
	tree.Add(childId1, NewTreeNode(RootUUID, "test1"))

	children, _  := tree.GetChildren(RootUUID)
	if len(children) == 1 && children[0] != childId1 {
		t.Errorf("Tree.GetChildren() returned %v, expected %v", children, []uuid.UUID{childId1})
	}


	tree.Remove(childId1) // Test removing a node

	children, _  = tree.GetChildren(RootUUID)
	if len(children) != 0 {
		t.Errorf("Tree.GetChildren() returned %v, expected %v", children, []uuid.UUID{})
	}
}

func TestTreeAncestor(t *testing.T) {
	tree := NewTree[string]()

	childId1 := uuid.New()
	tree.Add(childId1, NewTreeNode(RootUUID, "test1"))

	childId2 := uuid.New()
	tree.Add(childId2, NewTreeNode(childId1, "test2"))

	childId3 := uuid.New()
	tree.Add(childId3, NewTreeNode(childId2, "test3"))

	// Tree should be of the shape:
	// 	Root
	// 		|
	// 	   childId1
	// 		|
	// 	   childId2
	// 		|
	// 	   childId3
	
	// childId1 is an ancestor of childId3

	isAnc, _ := tree.IsAncestor(childId3, childId1)
	if !isAnc {
		t.Errorf("Tree.IsAncestor() returned false, expected true")
	}

	isRootAnc, _ := tree.IsAncestor(childId2, RootUUID)
	if !isRootAnc {
		t.Errorf("Tree.IsAncestor() returned false, expected true")
	}
}

func TestTreeDeleteSubTree(t *testing.T) {
	tree := NewTree[string]()

	childId1 := uuid.New()
	tree.Add(childId1, NewTreeNode(RootUUID, "test1"))

	childId2 := uuid.New()
	tree.Add(childId2, NewTreeNode(childId1, "test2"))

	childId3 := uuid.New()
	tree.Add(childId3, NewTreeNode(childId2, "test3"))

	// Tree should be of the shape:
	// 	Root
	// 		|
	// 	   childId1
	// 		|
	// 	   childId2
	// 		|
	// 	   childId3
	
	// Deleting subtree of childId1 should delete all nodes under it, and itself

	tree.DeleteSubTree(childId1)

	if tree.Contains(childId1) || tree.Contains(childId2) || tree.Contains(childId3) {
		t.Errorf("Tree.Contains() returned true, expected false")
	}
}

func TestTreeContains(t *testing.T) {
	tree := NewTree[string]()

	childId1 := uuid.New()
	tree.Add(childId1, NewTreeNode(RootUUID, "test1"))

	if !tree.Contains(childId1) {
		t.Errorf("Tree.Contains() returned false, expected true")
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