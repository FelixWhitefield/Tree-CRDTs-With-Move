package treecrdt

import (
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
)

type OpMove[MD Metadata,T opTimestamp[T]] struct {
	timestamp T
	parentID  uint64
	childID   uint64
	metadata  MD
}

type opTimestamp[T any] interface {
	clocks.TotalOrder[T]
	clocks.Timestamp[T]
}

func NewOpMove[MD Metadata, T opTimestamp[T]](timestamp T, parentID uint64, childID uint64, metadata MD) OpMove[MD, T] {
	return OpMove[MD, T]{timestamp: timestamp, parentID: parentID, childID: childID, metadata: metadata}
}

func (op OpMove[MD, T]) Timestamp() opTimestamp[T] {
	return op.timestamp
}

func (op OpMove[MD, T]) ParentID() uint64 {
	return op.parentID
}

func (op OpMove[MD, T]) ChildID() uint64 {
	return op.childID
}

func (op OpMove[MD, T]) Metadata() MD {
	return op.metadata
}

func (op OpMove[MD, T]) Compare(other OpMove[MD, T]) int {
	return op.timestamp.Compare(other.timestamp)
}
