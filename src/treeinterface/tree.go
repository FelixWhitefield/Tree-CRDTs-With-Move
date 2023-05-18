package treeinterface

import (
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

type Tree[MD any] interface {
	Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error) // Add a new node to the tree
	Delete(id uuid.UUID) error                                 // Delete a node from the tree
	Move(id uuid.UUID, newParentID uuid.UUID) error            // Move a node to a new parent
	Edit(id uuid.UUID, newMetadata MD) error                   // Edit the metadata of a node
	GetChildren(id uuid.UUID) ([]uuid.UUID, error)             // Get the children of a node
	GetParent(id uuid.UUID) (uuid.UUID, error)                 // Get the parent of a node
	Root() uuid.UUID                                           // Get the root of the tree
	GetMetadata(id uuid.UUID) (MD, error)                      // Get the metadata of a node
	Get(id uuid.UUID) (*tcrdt.TreeNode[MD], error)             // Get the node with the given ID
	GetTotalApplied() uint64                                   // Get the total number of operations applied to the tree
	ConnectionProvider() connection.ConnectionProvider         // Get the connection provider
}
