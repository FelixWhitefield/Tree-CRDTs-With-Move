package kleppmann

// Contains the CRDT state and implements the main algorithm
//
// `State` is independent of any peer, and should
// be equal between peers which have seen the same operations
//
// The op log is implemented as a linked list, as this allows for
// easy insertions in the middle of the list,
// as well as removals without the need to shift the list

import (
	"container/list"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"log"
)

type State[MD any, T opTimestamp[T]] struct {
	tree Tree[MD]   // state of the tree
	log  *list.List // ascending list of log moves
	// the log differs from the paper, as the paper states it should be descending
	// in practice this doesn't affect the algorithm
	extraConflict *TNConflict[MD]
}

func NewState[MD any, T opTimestamp[T]](conf *TNConflict[MD]) *State[MD, T] {
	return &State[MD, T]{tree: *NewTree[MD](), log: list.New(), extraConflict: conf}
}

// 'do_op' from the paper
// takes an op move, and applies it to the tree
// if the move is invalid, then the op is not applied but still logged
func (s *State[MD, T]) DoOp(op *OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP := s.tree.GetNode(op.ChldID)

	// If the child is an ancestor of the newParent
	isAnc, _ := s.tree.IsAncestor(op.NewP.PrntID, op.ChldID)
	newParentIsSelf := op.ChldID == op.NewP.PrntID

	// this is not in the algorithm.
	// it allows the user to define a custom conflict function
	conflict := false
	if s.extraConflict != nil {
		conflict = (*s.extraConflict)(op.NewP, &s.tree)
	}

	if !isAnc && !newParentIsSelf && !conflict && op.ChldID != s.tree.Root() {
		// instead of removing and then re-adding the node, we just move it
		// this ensures that the node will either be moved fully, or not at all
		// removing then adding may cause the node to be removed, but not added
		// which would cause the tree to be missing a node (which is unwanted)
		err := s.tree.Move(op.ChldID, op.NewP)

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
	s.tree.Remove(lop.op.ChldID)
	if lop.oldP != nil {
		s.tree.Add(lop.op.ChldID, lop.oldP)
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
func (s *State[MD, T]) ApplyOp(op *OpMove[MD, T]) {
	if s.log.Len() == 0 {
		logop := s.DoOp(op)
		s.log.PushBack(logop)
		return
	}
	e := s.log.Back()
	// This ignores the case where CompareOp returns 0, which is not defined in the paper
	// This should not happen in normal operation, if it does then the state is in an undefined state
	// loops while log op is greater than op
	for ; e != nil && e.Value.(*LogOpMove[MD, T]).CompareOp(op) == 1; e = e.Prev() {
		s.UndoOp(e.Value.(*LogOpMove[MD, T]))
	}

	// check if the op is already in the log (should not happen in normal operation)
	// or if we have moved to the front of the list
	if e == nil || !(e.Value.(*LogOpMove[MD, T]).CompareOp(op) == 0) {
		logop := s.DoOp(op)
		if e == nil { // If we have moved to the front of the list
			e = s.log.PushFront(logop)
		} else {
			e = s.log.InsertAfter(logop, e)
		}
	}
	e = e.Next()

	// redo ops until the end of the log
	for ; e != nil; e = e.Next() {
		s.RedoOp(e.Value.(*LogOpMove[MD, T]))
	}
}

// 'apply_ops' from the paper
// applies a list of ops to the tree
// the list of ops should be ordered, otherwise will be slow
func (s *State[MD, T]) ApplyOps(ops []*OpMove[MD, T]) {
	for _, op := range ops {
		s.ApplyOp(op)
	}
}

func (s *State[MD, T]) TruncateLogBefore(t T) {
	// oldest op is at the front of the list
	e := s.log.Front()
	for ; e != nil && e.Value.(*LogOpMove[MD, T]).op.Timestmp.Compare(t) == -1; e = e.Next() {
		s.log.Remove(e)
	}
}

func (s *State[MD, T]) Equals(other *State[MD, T]) bool {
	treeEq := s.tree.Equals(&other.tree)
	if !treeEq {
		return false
	}
	if s.log.Len() != other.log.Len() {
		return false
	}
	e1 := s.log.Front()
	e2 := other.log.Front()
	for ; e1 != nil && e2 != nil; e1, e2 = e1.Next(), e2.Next() {
		if e1.Value.(*LogOpMove[MD, T]).CompareOp(e2.Value.(*LogOpMove[MD, T]).op) != 0 {
			return false
		}
	}
	return true
}
