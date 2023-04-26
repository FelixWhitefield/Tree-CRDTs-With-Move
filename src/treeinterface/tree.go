package treeinterface

import (
	"github.com/google/uuid"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

type Tree[MD any] interface {
	Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	Move(id uuid.UUID, newParentID uuid.UUID) error
	Edit(id uuid.UUID, newMetadata MD) error
	GetChildren(id uuid.UUID) ([]uuid.UUID, error)
	GetParent(id uuid.UUID) (uuid.UUID, error)
	Root() uuid.UUID
	GetMetadata(id uuid.UUID) (MD, error)
	Get(id uuid.UUID) (*tcrdt.TreeNode[MD], error)
}