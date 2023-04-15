package treecrdt

type State[MD Metadata, T opTimestamp[T]] struct {
	tree Tree[MD]
	log []LogOpMove[MD, T]
}

func NewState[MD Metadata, T opTimestamp[T]]() State[MD, T] {
	return State[MD, T]{tree: *NewTree[MD](), log: []LogOpMove[MD, T]{}}
}

// 'do_op' from the paper
func (s *State[MD, T]) DoOp(op OpMove[MD, T]) LogOpMove[MD, T] {
	if s.tree.IsAncestor(op.childID, op.newParentID) || op.childID == op.newParentID {
		return NewLogOpMove(op, nil)
	} else {
		treeNode := NewTreeNode(op.newParentID, op.newMetadata)
		s.tree.Move(op.childID, *treeNode)
		return NewLogOpMove(op, treeNode)
	}
}

// 'undo_op' from the paper
// takes a log move, and moves the child back to its old parent
// if the old parent is nil, then the child is removed
func (s *State[MD, T]) UndoOp(lop LogOpMove[MD, T]) {
	if lop.oldP == nil {
		s.tree.Remove(lop.op.childID)
	} else {
		s.tree.Move(lop.op.childID, *lop.oldP)
	}
}

// 'redo_op' from the paper
// takes a log move, and applies the op move to the tree
func (s *State[MD, T]) RedoOp(lop LogOpMove[MD, T]) {
	op := lop.op
	logop := s.DoOp(op)
	s.log = append(s.log, logop)
}

// 'apply_op' from the paper
func (s *State[MD, T]) ApplyOp(op OpMove[MD, T]) {
	if (len(s.log) == 0) {
		s.DoOp(op)
	} else {
		for s.log[len(s.log)-1].CompareOp(op) == -1 {
			s.RedoOp(s.log[len(s.log)-1])
		}
	}
}

func (s *State[MD, T]) ApplyOps(ops []OpMove[MD, T]) {
	for _, op := range ops {
		s.ApplyOp(op)
	}
}