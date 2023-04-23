package k

// `TreeNode` is a node in the tree
// 
// `TreeNode` is a generic type, and so must be instantiated with a type for the metadata

import "github.com/google/uuid"

type TreeNode[MD any] struct {
	parentID uuid.UUID
	metadata MD
}

func NewTreeNode[MD any](parentID uuid.UUID, metadata MD) *TreeNode[MD] {
	return &TreeNode[MD]{parentID: parentID, metadata: metadata}
}

func (tn TreeNode[MD]) ParentID() uuid.UUID {
	return tn.parentID
}

func (tn TreeNode[MD]) Metadata() MD {
	return tn.metadata
}

// defines a conflict (or multiple conflicts)
// the conflict will be between the tree node to be inserted and the current state of the tree
// if the node were to cause a conflict, the function should return true
// the function should not modify the tree or the tree node
type TNConflict[MD any] func(tn1 *TreeNode[MD], tn2 *Tree[MD]) bool
