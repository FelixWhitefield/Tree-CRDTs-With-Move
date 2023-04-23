package treeinterface

import (
	"github.com/google/uuid"
)

type Tree[MD any] interface {
	Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	Move(id uuid.UUID, newParentID uuid.UUID) error
	Edit(id uuid.UUID, newMetadata MD) error
	GetChildren(id uuid.UUID) ([]uuid.UUID, error)
	GetParent(id uuid.UUID) (uuid.UUID, error)
	Root() uuid.UUID
	Get(id uuid.UUID) (MD, error)
}