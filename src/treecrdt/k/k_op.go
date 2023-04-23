package k
// `OpMove` is a move operation
//
// Contains method to compare two operations 

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/google/uuid"
)

// Could use the following to distinguish between different types of operations (add, remove, move)
// Could be used for implementing optimisations (as discuessed in the paper)
// Add Op: Has all childID, parentID, and metadata (childID is not in tree)
// Remove Op: Has all childID and nil parentID
// Move Op: Has all childID, parentID and metadata

// Represents moving node with id childID to parent and metadata within newP
type OpMove[MD any, T opTimestamp[T]] struct {
	timestamp T
	childID   uuid.UUID
	newP      *TreeNode[MD]
}

type opTimestamp[T any] interface {
	c.TotalOrder[T]
	c.Timestamp[T]
	ActorID() uuid.UUID
}

func NewOpMove[MD any, T opTimestamp[T]](timestamp T, parentID uuid.UUID, childID uuid.UUID, metadata MD) *OpMove[MD, T] {
	return &OpMove[MD, T]{timestamp: timestamp, childID: childID, newP: NewTreeNode(parentID, metadata)}
}

func (op OpMove[MD, T]) Timestamp() T {
	return op.timestamp
}

func (op OpMove[MD, T]) ParentID() uuid.UUID {
	return op.newP.parentID
}

func (op OpMove[MD, T]) ChildID() uuid.UUID {
	return op.childID
}

func (op OpMove[MD, T]) Metadata() MD {
	return op.newP.metadata
}

// Compares two OpMoves by their timestamps
func (op *OpMove[MD, T]) Compare(other *OpMove[MD, T]) int {
	return op.timestamp.Compare(other.timestamp)
}
