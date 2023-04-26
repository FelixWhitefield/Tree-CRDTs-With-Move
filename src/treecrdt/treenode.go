package treecrdt

// `TreeNode` is a node in the tree
//
// `TreeNode` is a generic type, and so must be instantiated with a type for the metadata

import "github.com/google/uuid"

type TreeNode[MD any] struct {
	PrntID uuid.UUID
	Meta   MD
}

func NewTreeNode[MD any](parentID uuid.UUID, metadata MD) *TreeNode[MD] {
	return &TreeNode[MD]{PrntID: parentID, Meta: metadata}
}

func (tn TreeNode[MD]) ParentID() uuid.UUID {
	return tn.PrntID
}

func (tn TreeNode[MD]) Metadata() MD {
	return tn.Meta
}

// defines a conflict (or multiple conflicts)
// the conflict will be between the tree node to be inserted and the current state of the tree
// if the node were to cause a conflict, the function should return true
// the function should not modify the tree or the tree node
type TNConflict[MD any] func(tn1 *TreeNode[MD], tn2 *Tree[MD]) bool

func (tn *TreeNode[MD]) Equals(other *TreeNode[MD]) bool {
	return tn.PrntID == other.PrntID
}
