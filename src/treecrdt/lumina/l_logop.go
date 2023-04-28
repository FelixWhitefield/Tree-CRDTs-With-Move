package lumina

import (
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

// Represents a log of the move `op`, and the old parent `oldP`
type LogOpMove[MD any, T opTimestamp[T]] struct {
	op   *OpMove[MD, T]
	oldP *TreeNode[MD]
}

func NewLogOpMove[MD any, T opTimestamp[T]](op *OpMove[MD, T], oldP *TreeNode[MD]) *LogOpMove[MD, T] {
	return &LogOpMove[MD, T]{op: op, oldP: oldP}
}

func (lop *LogOpMove[MD, T]) Timestamp() opTimestamp[T] {
	return lop.op.Timestmp
}

func (lop *LogOpMove[MD, T]) OpMove() OpMove[MD, T] {
	return *lop.op
}

// Compares a LogOpMove with an OpMove
// This is useful for the state
func (lop *LogOpMove[MD, T]) CompareOp(other *OpMove[MD, T]) int {
	return lop.op.Timestmp.Compare(other.Timestmp)
}

// Compares a LogOpMove's priority with an OpMove's priority
// This is used for determining the order of concurrent operations
func (lop *LogOpMove[MD, T]) ComparePriority(other *OpMove[MD, T]) int {
	return lop.op.Priotity.Compare(&other.Priotity)
}

// Checks if two LogOpMoves are equal
func (lop *LogOpMove[MD, T]) Equals(other *LogOpMove[MD, T]) bool {
	return lop.op.Timestmp.Compare(other.op.Timestmp) == 0 && lop.op.Priotity.Compare(&other.op.Priotity) == 0
}
