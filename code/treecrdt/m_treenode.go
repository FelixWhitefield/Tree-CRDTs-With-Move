package treecrdt

import "github.com/google/uuid"

type TreeNode[MD Metadata] struct {
	parentID uuid.UUID
	metadata MD
}

func NewTreeNode[MD Metadata](parentID uuid.UUID, metadata MD) *TreeNode[MD] {
	return &TreeNode[MD]{parentID: parentID, metadata: metadata}
}

func (tn TreeNode[MD]) ParentID() uuid.UUID {
	return tn.parentID
}

func (tn TreeNode[MD]) Metadata() MD {
	return tn.metadata
}

type TNConflict[MD Metadata] func(tn1 TreeNode[MD], tn2 TreeNode[MD]) bool