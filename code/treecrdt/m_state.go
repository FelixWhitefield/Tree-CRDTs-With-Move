package treecrdt

import (
	"container/list"
)

type State[MD Metadata, T opTimestamp[T]] struct {
	tree Tree[MD]
	log  list.List
}

func NewState[MD Metadata, T opTimestamp[T]]() State[MD, T] {
	return State[MD, T]{tree: *NewTree[MD](), log: *list.New()}
}

// 'do_op' from the paper
// takes an op move, and applies it to the tree
// if the move is invalid, then the op is not applied but still logged
func (s *State[MD, T]) DoOp(op OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP, _ := s.tree.GetNode(op.childID)
	if !(s.tree.IsAncestor(op.childID, op.newParentID) || op.childID == op.newParentID) {
		s.tree.Move(op.childID, oldP)
	}
	return NewLogOpMove(op, oldP)
}

// 'undo_op' from the paper
// takes a log move, and moves the child back to its old parent
// if the old parent is nil, then the child is removed
func (s *State[MD, T]) UndoOp(lop *LogOpMove[MD, T]) {
	if lop.oldP == nil {
		s.tree.Remove(lop.op.childID)
	} else {
		s.tree.Move(lop.op.childID, lop.oldP)
	}
}

// 'redo_op' from the paper
// takes a log move, and applies the op move to the tree
func (s *State[MD, T]) RedoOp(lop *LogOpMove[MD, T]) {
	op := lop.op
	logop := s.DoOp(op)
	*lop = *logop
}

// 'apply_op' from the paper
// applies an op to the tree
// undo's and redo's ops if necessary
func (s *State[MD, T]) ApplyOp(op OpMove[MD, T]) {
	if s.log.Len() == 0 {
		logop := s.DoOp(op)
		s.log.PushBack(logop)
	} else {
		e := s.log.Back()
		for ; e.Value.(LogOpMove[MD, T]).CompareOp(op) == -1; e = e.Prev() {
			s.UndoOp(e.Value.(*LogOpMove[MD, T]))
		}
		logop := s.DoOp(op)
		newLogElem := s.log.InsertAfter(logop, e)

		e = newLogElem.Next()
		for ; e != nil; e = e.Next() {
			s.RedoOp(e.Value.(*LogOpMove[MD, T]))
		}
	}
}

// 'apply_ops' from the paper
// applies a list of ops to the tree
// the list of ops should be ordered, otherwise will be slow
func (s *State[MD, T]) ApplyOps(ops []OpMove[MD, T]) {
	for _, op := range ops {
		s.ApplyOp(op)
	}
}
