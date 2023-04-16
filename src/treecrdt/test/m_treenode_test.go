package treecrdt_test

import (
	"testing"
	"github.com/google/uuid"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

func TestTreeNode(t *testing.T) {
	u1 := uuid.New()
	tn := NewTreeNode[string](u1, "test")
	if tn == nil {
		t.Errorf("NewTreeNode() returned nil")
	}
	if tn.ParentID() != u1 {
		t.Errorf("NewTreeNode() returned wrong parentID")
	}
	if tn.Metadata() != "test" {
		t.Errorf("NewTreeNode() returned wrong metadata")
	}
}

func TestTreeNodeConflict(t *testing.T) {
	tnc := TNConflict[Metadata](func(tn1 TreeNode[Metadata], tn2 TreeNode[Metadata]) bool {
		return tn1.Metadata() == tn2.Metadata()
	})

	tn1 := NewTreeNode[Metadata](uuid.New(), "test")
	tn2 := NewTreeNode[Metadata](uuid.New(), "test")
	tn3 := NewTreeNode[Metadata](uuid.New(), "test2")
	if !tnc(*tn1, *tn2) {
		t.Errorf("TNConflict() returned false")
	}
	if tnc(*tn1, *tn3) {
		t.Errorf("TNConflict() returned true")
	}
}