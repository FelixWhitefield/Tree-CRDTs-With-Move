package treecrdt

type LogOpMove[MD Metadata, T opTimestamp[T]] struct {
	op   OpMove[MD, T]
	oldP *TreeNode[MD]
}

func NewLogOpMove[MD Metadata, T opTimestamp[T]](op OpMove[MD, T], oldP *TreeNode[MD]) *LogOpMove[MD, T] {
	return &LogOpMove[MD, T]{op: op, oldP: oldP}
}

func (lop LogOpMove[MD, T]) Timestamp() opTimestamp[T] {
	return lop.op.timestamp
}

func (lop LogOpMove[MD, T]) OpMove() OpMove[MD, T] {
	return lop.op
}

func (lop LogOpMove[MD, T]) CompareOp(other OpMove[MD, T]) int {
	return lop.op.timestamp.Compare(other.timestamp)
}
