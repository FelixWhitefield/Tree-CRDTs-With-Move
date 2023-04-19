package treecrdt_test

import (
	"testing"
	"github.com/google/uuid"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

func TestTreeNode(t *testing.T) {
	testData := struct{
		parentID uuid.UUID
		metadata string
	}{
		uuid.New(),
		"test",
	}

	tn := NewTreeNode(testData.parentID, testData.metadata)
	if tn == nil {
		t.Errorf("NewTreeNode() returned nil")
	}
	if tn != nil && tn.ParentID() != testData.parentID {
		t.Errorf("NewTreeNode() returned wrong parentID")
	}
	if tn != nil && tn.Metadata() != testData.metadata {
		t.Errorf("NewTreeNode() returned wrong metadata")
	}
}

func TestTreeNodeConflict(t *testing.T) {
	tnc := TNConflict[Metadata](func(tn *TreeNode[Metadata], t *Tree[Metadata]) bool {
		children, ok := t.GetChildren(tn.ParentID())
		if !ok {
			return false
		}
		for _, tnid := range children {
			if t.GetNode(tnid).Metadata() == tn.Metadata() {
				return true
			}
		}
		return false;
	})

	tn1 := NewTreeNode[Metadata](RootUUID, "test")
	tn2 := NewTreeNode[Metadata](RootUUID, "test")

	tree := NewTree[Metadata]()
	u1 := uuid.New()
	tree.Add(u1, tn1)

	if !tnc(tn2, tree) {
		t.Errorf("tnc() returned false")
	}
}