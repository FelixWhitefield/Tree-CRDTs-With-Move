package treecrdt

type OldData[MD Metadata] struct {
	parentID uint64
	metadata MD
}

type LogOpMove[MD Metadata, T opTimestamp[T]] struct {
	op   OpMove[MD, T]
	oldP OldData[MD]
}

func NewLogOpMove[MD Metadata, T opTimestamp[T]](op OpMove[MD, T], oldP OldData[MD]) LogOpMove[MD, T] {
	return LogOpMove[MD, T]{op: op, oldP: oldP}
}

func (lop LogOpMove[MD, T]) Timestamp() opTimestamp[T] {
	return lop.op.timestamp
}

func (lop LogOpMove[MD, T]) Compare(other LogOpMove[MD, T]) int {
	return lop.op.timestamp.Compare(other.op.timestamp)
}
