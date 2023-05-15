package lumina

// Represents the state of the CRDT

import (
	"container/list"
	//"fmt"
	//"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

// Represents the state of the CRDT
type State[MD any, T opTimestamp[T]] struct {
	tree    Tree[MD]
	moveLog *list.List // List of move operations
}

func NewLState[MD any, T opTimestamp[T]]() *State[MD, T] {
	return &State[MD, T]{tree: *NewTree[MD](), moveLog: list.New()}
}

// Applies an operation to the state
func (s *State[MD, T]) ApplyOp(op Operation[T]) {
	switch op := op.(type) {
	case *OpMove[MD, T]:
		s.ApplyMoveOp(op)
	case *OpAdd[MD, T]:
		parentInTree := s.tree.WithinTree(op.NewP.ParentID()) != nil
		childInTree := s.tree.WithinTree(op.ChldID) != nil // If child is in the tree

		if !childInTree && parentInTree { // If child is not in the tree and parent is in the tree
			s.tree.Add(op.ChldID, op.NewP)
		}
	case *OpRemove[T]:
		if op.ChldID != s.tree.Root() { // If child is not the root
			withinTree := s.tree.WithinTree(op.ChldID) != nil
			if withinTree {
				s.tree.Remove(op.ChldID)
			}
		}
	}
	//fmt.Println(s.tree.String())
}

// Applies a move operation to the tree
// This differs from maram, as it does not implement the conflict
// resolution policy the same way.
// This instead uses the priority to determine a total order of concurrent operations
// and undo's and redo's operations accordingly (Similar to Kleppmann's algorithm)
func (s *State[MD, T]) ApplyMoveOp(opMov *OpMove[MD, T]) {
	if s.moveLog.Len() == 0 {
		logop := s.DoMoveOp(opMov)
		s.moveLog.PushBack(logop)
		return
	}

	e := s.moveLog.Back()

	for ; e != nil && e.Value.(*LogOpMove[MD, T]).CompareOp(opMov) == 2 && e.Value.(*LogOpMove[MD, T]).ComparePriority(opMov) == 1; e = e.Prev() {
		s.UndoMoveOp(e.Value.(*LogOpMove[MD, T]))
	}

	if e == nil || !(e.Value.(*LogOpMove[MD, T]).ComparePriority(opMov) == 0) { // Got to the front of the list and not in tree
		logop := s.DoMoveOp(opMov)
		if e == nil {
			e = s.moveLog.PushFront(logop)
		} else {
			e = s.moveLog.InsertAfter(logop, e)
		}
	}

	e = e.Next()

	// redo ops until the end of the log
	for ; e != nil; e = e.Next() {
		s.RedoMoveOp(e.Value.(*LogOpMove[MD, T]))
	}
}

// Undo's the log operation
func (s *State[MD, T]) UndoMoveOp(lop *LogOpMove[MD, T]) {
	s.tree.Move(lop.op.ChldID, lop.oldP)
}

// Redo's the log operation
func (s *State[MD, T]) RedoMoveOp(lop *LogOpMove[MD, T]) {
	logop := s.DoMoveOp(lop.op)
	*lop = *logop
}

// Does the move operation
// This checks if the move operation is valid and then executes it
func (s *State[MD, T]) DoMoveOp(opMov *OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP := s.tree.WithinTree(opMov.ChldID).Node
	childIsRoot := opMov.ChldID == s.tree.Root()
	//parentInTree := s.tree.WithinTree(opMov.NewP.PrntID) != nil // If parent is in the tree
	//childInTree := s.tree.WithinTree(opMov.ChldID) != nil                  // If child is in the tree (may be child of tombstone)
	childIsAnc, _ := s.tree.IsAncestor(opMov.NewP.PrntID, opMov.ChldID) // If parent is an ancestor of child

	if !childIsAnc && !childIsRoot && opMov.ChldID != opMov.NewP.PrntID {
		s.tree.Move(opMov.ChldID, opMov.NewP)
	}
	return NewLogOpMove(opMov, oldP)
}

// Checks if the state is equal to another state
func (s *State[MD, T]) Equals(other *State[MD, T]) bool {
	treeEq := s.tree.Equals(&other.tree)
	return treeEq
}
