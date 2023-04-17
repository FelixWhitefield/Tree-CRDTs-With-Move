package treecrdt_test

// import (
// 	"testing"
// 	"github.com/google/uuid"
// 	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
// )

// func TestTreeNode(t *testing.T) {
// 	u1 := uuid.New()
// 	tn := NewTreeNode[string](u1, "test")
// 	if tn == nil {
// 		t.Errorf("NewTreeNode() returned nil")
// 	}
// 	if tn.ParentID() != u1 {
// 		t.Errorf("NewTreeNode() returned wrong parentID")
// 	}
// 	if tn.Metadata() != "test" {
// 		t.Errorf("NewTreeNode() returned wrong metadata")
// 	}
// }

// func TestTreeNodeConflict(t *testing.T) {
// 	tnc := TNConflict[Metadata](func(tn *TreeNode[Metadata], t *Tree[Metadata]) bool {
// 		children, ok := t.GetChildren(tn.ParentID())
// 		if !ok {
// 			return false
// 		}
// 		for _, tnid := range children {
// 			if t.GetNode(tnid).Metadata() == tn.Metadata() {
// 				return true
// 			}
// 		}
// 	})

// 	tn1 := NewTreeNode[Metadata](uuid.New(), "test")
// 	tn2 := NewTreeNode[Metadata](uuid.New(), "test")
// 	tn3 := NewTreeNode[Metadata](uuid.New(), "test2")

// }