package lumina

// This was the original state implementation, but it was replaced with the
// new state implementation in src\treecrdt\lumina\l_state.go
// This is kept here for reference
//
// The code here does not work as intended, so a new state and tree implementation
// was created to fix the issues.

import (
	"container/list"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

// Represents the state of the CRDT
type State[MD any, T opTimestamp[T]] struct {
	tree    treecrdt.Tree[MD]
	moveLog *list.List // List of move operations
}

func NewState[MD any, T opTimestamp[T]]() *State[MD, T] {
	return &State[MD, T]{tree: *treecrdt.NewTree[MD](), moveLog: list.New()}
}

// Applies an operation to the state
func (s *State[MD, T]) ApplyOp(op Operation[T]) {
	switch op := op.(type) {
	case *OpMove[MD, T]:
		s.ApplyMoveOp(op)
	case *OpAdd[MD, T]:

		parentInTree := !s.tree.IsDeleted(op.NewP.ParentID())
		childInTree, _ := s.tree.IsAncestor(op.ChldID, s.tree.Root()) // If child is in the tree
		if !childInTree && parentInTree {                             // If child is not in the tree and parent is in the tree
			s.tree.Add(op.ChldID, op.NewP)
		}
	case *OpRemove[T]:
		if op.ChldID != s.tree.Root() { // If child is not the root
			node := s.tree.GetNode(op.ChldID)
			if node != nil {
				s.tree.Move(op.ChldID, treecrdt.NewTreeNode(s.tree.Tombstone(), node.Meta)) // Move child to the tombstone
			}
		}
	}
}

// Applies a move operation to the tree
// This differs from maram, as it does not implement the conflict
// resolution policy the same way.
// This instead uses the priority to determine a total order of concurrent operations
// and undo's and redo's operations accordingly (Similar to Kleppmann's algorithm)
func (s *State[MD, T]) ApplyMoveOp(opMov *OpMove[MD, T]) {
	if s.moveLog.Len() == 0 {
		logop := s.DoMoveOp(opMov)
		if logop != nil {
			s.moveLog.PushBack(logop)
		}
		return
	}
	e := s.moveLog.Back()

	for e != nil && e.Value.(*LogOpMove[MD, T]).CompareOp(opMov) == 2 && e.Value.(*LogOpMove[MD, T]).ComparePriority(opMov) == 1 {
		if s.UndoMoveOp(e.Value.(*LogOpMove[MD, T])) {
			e = e.Prev()
		} else {
			cur := e
			e = e.Prev()
			s.moveLog.Remove(cur)
		}
	}

	if e == nil || !(e.Value.(*LogOpMove[MD, T]).ComparePriority(opMov) == 0) { // Got to the front of the list and not in tree
		logop := s.DoMoveOp(opMov)
		if logop != nil {
			if e == nil {
				e = s.moveLog.PushFront(logop)
			} else {
				e = s.moveLog.InsertAfter(logop, e)
			}
			e = e.Next()
		}
	}

	// redo ops until the end of the log
	for e != nil {
		logop := s.RedoMoveOp(e.Value.(*LogOpMove[MD, T]))
		if logop == nil {
			cur := e
			e = e.Next()
			s.moveLog.Remove(cur)
		} else {
			e = e.Next()
		}
	}
}

// Undo's the log operation
func (s *State[MD, T]) UndoMoveOp(lop *LogOpMove[MD, T]) bool {
	if s.tree.IsDeleted(lop.op.ChldID) {
		return false
	}
	s.tree.Remove(lop.op.ChldID)
	s.tree.Add(lop.op.ChldID, lop.oldP)
	return true
}

// Redo's the log operation
func (s *State[MD, T]) RedoMoveOp(lop *LogOpMove[MD, T]) *LogOpMove[MD, T] {
	if s.tree.IsDeleted(lop.op.ChldID) {
		return nil
	}
	logop := s.DoMoveOp(lop.op)
	return logop
}

// Does the move operation
// This checks if the move operation is valid and then executes it
func (s *State[MD, T]) DoMoveOp(opMov *OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP := s.tree.GetNode(opMov.ChldID)
	if oldP == nil {
		return nil
	}
	childIsRoot := opMov.ChldID == s.tree.Root()
	parentInTree := s.tree.Anywhere(opMov.NewP.ParentID()) // If parent is in the tree
	if opMov.NewP.ParentID() == s.tree.Root() {
		parentInTree = true
	}
	childInTree := s.tree.TotalContains(opMov.ChldID)                       // If child is in the tree (may be child of tombstone)
	childIsAnc, _ := s.tree.IsAncestor(opMov.NewP.ParentID(), opMov.ChldID) // If parent is an ancestor of child
	if !childIsAnc && parentInTree && childInTree && !childIsRoot {
		s.tree.Move(opMov.ChldID, opMov.NewP)
		return NewLogOpMove(opMov, oldP)
	}
	return nil
}

// Checks if the state is equal to another state
func (s *State[MD, T]) Equals(other *State[MD, T]) bool {
	treeEq := s.tree.Equals(&other.tree)
	return treeEq
}
