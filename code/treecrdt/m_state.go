package treecrdt

type State[MD Metadata, T opTimestamp[T]] struct {
	tree Tree[MD]
	log []LogOpMove[MD, T]
}

func NewState[MD Metadata, T opTimestamp[T]]() State[MD, T] {
	return State[MD, T]{tree: *NewTree[MD](), log: []LogOpMove[MD, T]{}}
}

func (s *State[MD, T]) DoOp(op OpMove[MD, T]) {
	//
}

func (s *State[MD, T]) UndoOp(op OpMove[MD, T]) {
	//
}

func (s *State[MD, T]) RedoOp(op OpMove[MD, T]) {
	//
}

func (s *State[MD, T]) ApplyOp(op OpMove[MD, T]) {
	//
}

func (s *State[MD, T]) ApplyOps(ops []OpMove[MD, T]) {
	//
}