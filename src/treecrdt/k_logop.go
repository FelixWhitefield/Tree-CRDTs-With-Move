package treecrdt

// Represents a log of the move `op`, and the old parent `oldP`
type LogOpMove[MD Metadata, T opTimestamp[T]] struct {
	op   *OpMove[MD, T]
	oldP *TreeNode[MD]
}

func NewLogOpMove[MD Metadata, T opTimestamp[T]](op *OpMove[MD, T], oldP *TreeNode[MD]) *LogOpMove[MD, T] {
	return &LogOpMove[MD, T]{op: op, oldP: oldP}
}

func (lop LogOpMove[MD, T]) Timestamp() opTimestamp[T] {
	return lop.op.timestamp
}

func (lop LogOpMove[MD, T]) OpMove() OpMove[MD, T] {
	return *lop.op
}

// Compares a LogOpMove with an OpMove
// This is useful for the state
func (lop LogOpMove[MD, T]) CompareOp(other *OpMove[MD, T]) int {
	return lop.op.timestamp.Compare(other.timestamp)
}
