package lumina

import (
	"container/list"

	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)

// Represents the state of the CRDT
type State[MD any, T opTimestamp[T]] struct {
	tree    Tree[MD]
	moveLog *list.List // List of move operations
}

func NewState[MD any, T opTimestamp[T]]() *State[MD, T] {
	return &State[MD, T]{tree: *NewTree[MD](), moveLog: list.New()}
}

// Applies an operation to the state
func (s *State[MD, T]) ApplyOp(op Operation[T]) {
	switch op := op.(type) {
	case *OpMove[MD, T]:
		s.ApplyMoveOp(op)
	case *OpAdd[MD, T]:
		parentInTree, _ := s.tree.IsAncestor(op.NewP.ParentID(), s.tree.Root()) // If parent is in the tree
		if op.NewP.ParentID() == s.tree.Root() {
			parentInTree = true
		}
		childInTree, _ := s.tree.IsAncestor(op.ChldID, s.tree.Root()) // If child is in the tree
		if !childInTree && parentInTree {                             // If child is not in the tree and parent is in the tree
			s.tree.Add(op.ChldID, op.NewP)
		}
	case *OpRemove[T]:
		if op.ChldID != s.tree.Root() { // If child is not the root
			node := s.tree.GetNode(op.ChldID)
			if node != nil {
				s.tree.Move(op.ChldID, NewTreeNode(s.tree.Root(), node.Meta)) // Move child to the tombstone
			}
		}
	}
}

// Applies a move operation to the tree
// This differs from the algorithm in the paper, as it does not implement the conflict
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
	s.tree.Remove(lop.op.ChldID)
	if lop.oldP != nil {
		s.tree.Add(lop.op.ChldID, lop.oldP)
	}
}

// Redo's the log operation
func (s *State[MD, T]) RedoMoveOp(lopMov *LogOpMove[MD, T]) {
	logop := s.DoMoveOp(lopMov.op)
	*lopMov = *logop
}

// Does the move operation
// This checks if the move operation is valid and then executes it
func (s *State[MD, T]) DoMoveOp(opMov *OpMove[MD, T]) *LogOpMove[MD, T] {
	oldP := s.tree.GetNode(opMov.ChldID)
	childIsRoot := opMov.ChldID == s.tree.Root()
	parentInTree, _ := s.tree.IsAncestor(opMov.NewP.ParentID(), s.tree.Root()) // If parent is in the tree
	if opMov.NewP.ParentID() == s.tree.Root() {
		parentInTree = true
	}
	childInTree := s.tree.TotalContains(opMov.ChldID) // If child is in the tree (may be child of tombstone)
	childRemoved := s.tree.GetNode(opMov.ChldID).PrntID == s.tree.Tombstone()
	childIsAnc, _ := s.tree.IsAncestor(opMov.NewP.ParentID(), opMov.ChldID) // If parent is an ancestor of child
	if !childIsAnc && parentInTree && childInTree && !childIsRoot && !childRemoved {
		s.tree.Move(opMov.ChldID, opMov.NewP)
	}
	return NewLogOpMove(opMov, oldP)
}

// Checks if the state is equal to another state
func (s *State[MD, T]) Equals(other *State[MD, T]) bool {
	treeEq := s.tree.Equals(&other.tree)
	if !treeEq {
		return false
	}
	if s.moveLog.Len() != other.moveLog.Len() {
		return false
	}
	e1 := s.moveLog.Front()
	e2 := other.moveLog.Front()
	for ; e1 != nil && e2 != nil; e1, e2 = e1.Next(), e2.Next() {
		if !e1.Value.(*LogOpMove[MD, T]).Equals(e2.Value.(*LogOpMove[MD, T])) {
			return false
		}
	}
	return true
}

// Attempted implementation of algorithm from the paper
// func (s *State[MD, T]) ApplyMoveOp(opMov *OpMove[MD, T]) {
// 	if s.moveLog.Len() == 0 {
// 		s.tree.Move(opMov.ChldID, opMov.NewP)
// 		s.moveLog.PushBack(opMov.ChldID)
// 		return
// 	}

// 	e := s.moveLog.Back()
// 	atLeastOne := false

// 	for ; e != nil && e.Value.(*LogOpMove[MD, T]).CompareOp(opMov) == 2; e = e.Prev() {
// 		atLeastOne = true
// 		s.UndoMoveOp(e.Value.(*LogOpMove[MD, T]))
// 	}

// 	if e == nil { // Got to the front of the list
// 		e = s.moveLog.Front()
// 	} else if e.Value.(*OpMove[MD, T]).CompareOp(opMov) == 0 { // Value already in tree
// 		return
// 	} else { // Move to the next element (which is concurrent)
// 		e = e.Next()
// 	}

// 	if atLeastOne {
// 		for ; e != nil; e = e.Next() {
// 			logOpMove := e.Value.(*LogOpMove[MD, T])

// 			// If child is the same
// 			if logOpMove.op.ChldID == opMov.ChldID {
// 				if opMov.Priotity.Compare(&logOpMove.op.Priotity) == 1 {
// 					s.DoMoveOp(opMov)
// 					continue
// 				} else {
// 					s.DoLogMoveOp(logOpMove)
// 				}
// 			}

// 		}
// 	} else {
// 		s.DoMoveOp(opMov)
// 	}
// }
