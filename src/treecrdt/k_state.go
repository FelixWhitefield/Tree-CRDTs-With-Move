package treecrdt

// Contains the CRDT state and implements the algorithm
//
// `State` is independent of any peer, and should 
// be equal between peers which have seens the same operations
//
// The code differes slightly from the paper
// However, the algorithm functions the same

import (
	"container/list"
	"log"
)

type State[MD Metadata, T opTimestamp[T]] struct {
	tree Tree[MD] // state of the tree
	log  list.List // ascending list of log moves 
	// the log differs from the paper, as the paper states it should be descending
	// in practice this doesn't affect the algorithm
	extraConflict *TNConflict[MD] 
}

func NewState[MD Metadata, T opTimestamp[T]]() State[MD, T] {
	return State[MD, T]{tree: *NewTree[MD](), log: *list.New()}
}

// 'do_op' from the paper
// takes an op move, and applies it to the tree
// if the move is invalid, then the op is not applied but still logged
func (s *State[MD, T]) DoOp(op OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP := s.tree.GetNode(op.childID)

	isAnc, _ := s.tree.IsAncestor(op.childID, op.newP.parentID)
	newParentIsSelf := op.childID == op.newP.parentID

	// this is not in the algorithm.
	// it allows the user to define a custom conflict function
	conflict := false
	if s.extraConflict != nil && (*s.extraConflict)(op.newP, &s.tree) {
		conflict = true
	}

	if !isAnc && !newParentIsSelf && !conflict {
		// instead of removing and then re-adding the node, we just move it
		// this ensures that the node will either be moved fully, or not at all
		// removing then adding may cause the node to be removed, but not added
		// which would cause the tree to be missing a node (which is unwanted)
		err := s.tree.Move(op.childID, op.newP)

		// errors will happen during concurrent operations
		// this is normal, and will be resolved once the operations are applied in order
		if err != nil {
			log.Println("Error moving node: ", err)
		}
	}

	return NewLogOpMove(op, oldP)
}

// 'undo_op' from the paper
// takes a log move, and moves the child back to its old parent
// if the old parent is nil, then the child is removed
func (s *State[MD, T]) UndoOp(lop *LogOpMove[MD, T]) {
	s.tree.Remove(lop.op.childID)
	if !(lop.oldP == nil) {
		s.tree.Add(lop.op.childID, lop.oldP)
	}
}

// 'redo_op' from the paper
// takes a log move, and applies the op move to the tree
func (s *State[MD, T]) RedoOp(lop *LogOpMove[MD, T]) {
	logop := s.DoOp(lop.op)
	*lop = *logop // update the logop in place (optimisation)
}

// 'apply_op' from the paper
// applies an op to the tree
// undo's and redo's ops if necessary
// the paper defines this method as a recursive function
// this implementation is iterative, and does not remove and then re-add the op to the logas - this would be inefficient
// instead, a linked list is used to store the log, and the op is inserted in the correct place
// the elements in the list are modified in place
func (s *State[MD, T]) ApplyOp(op OpMove[MD, T]) {
	if s.log.Len() == 0 {
		logop := s.DoOp(op)
		s.log.PushBack(logop)
	} else {
		e := s.log.Back()
		// This ignores the case where CompareOp returns 0, which is not defined in the paper
		// This should not happen in normal operation, if it does then the state is in an undefined state
		// loops while log op is greater than op
		for ; e.Value.(LogOpMove[MD, T]).CompareOp(op) == 1; e = e.Prev() {
			s.UndoOp(e.Value.(*LogOpMove[MD, T]))
		}

		// check if the op is already in the log (should not happen in normal operation)
		if !(e.Value.(LogOpMove[MD, T]).CompareOp(op) == 0) {
			logop := s.DoOp(op)
			e = s.log.InsertAfter(logop, e)
		}
		e = e.Next()
		
		// redo ops until the end of the log
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