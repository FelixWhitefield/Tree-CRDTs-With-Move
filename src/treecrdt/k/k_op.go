package k

// `OpMove` is a move operation
//
// Contains method to compare two operations

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

// Could use the following to distinguish between different types of operations (add, remove, move)
// Could be used for implementing optimisations (as discuessed in the paper)
// Add Op: Has all childID, parentID, and metadata (childID is not in tree)
// Remove Op: Has all childID and nil parentID
// Move Op: Has all childID, parentID and metadata

// Represents moving node with id childID to parent and metadata within newP
type OpMove[MD any, T opTimestamp[T]] struct {
	Timestmp T
	ChldID   uuid.UUID
	NewP     *TreeNode[MD]
}

type opTimestamp[T any] interface {
	c.TotalOrder[T]
	c.Timestamp[T]
	ActorID() uuid.UUID
}

func NewOpMove[MD any, T opTimestamp[T]](timestamp T, parentID uuid.UUID, childID uuid.UUID, metadata MD) *OpMove[MD, T] {
	return &OpMove[MD, T]{Timestmp: timestamp, ChldID: childID, NewP: NewTreeNode(parentID, metadata)}
}

func (op OpMove[MD, T]) Timestamp() T {
	return op.Timestmp
}

func (op OpMove[MD, T]) ParentID() uuid.UUID {
	return op.NewP.PrntID
}

func (op OpMove[MD, T]) ChildID() uuid.UUID {
	return op.ChldID
}

func (op OpMove[MD, T]) Metadata() MD {
	return op.NewP.Meta
}

// Compares two OpMoves by their timestamps
func (op *OpMove[MD, T]) Compare(other *OpMove[MD, T]) int {
	return op.Timestmp.Compare(other.Timestmp)
}
