package realmaram

import (
	"container/list"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

// Represents the state of the CRDT
type State[MD any, T opTimestamp[T]] struct {
	tree    Tree[MD]
	moveOps *list.List //[]OpMove[MD, T] // List of move operations
	ranks   map[uuid.UUID]int // Map of node ids to their rank
}

func NewState[MD any, T opTimestamp[T]]() *State[MD, T] {
	return &State[MD, T]{tree: *NewTree[MD](), moveOps: list.New(), ranks: make(map[uuid.UUID]int)}
}

// Applies an operation to the state
func (s *State[MD, T]) ApplyOp(op Operation[T]) {
	switch op := op.(type) {
	case OpMove[MD, T]:
		s.ApplyMoveOp(op)
	case *OpAdd[MD, T]:
		node := s.tree.GetNode(op.ChldID)
		// Checks: requires { [@expl:add1] not F.mem n s.nodes }
		if node != nil {
			return
		}
		// Checks: requires { [@expl:add2] F.mem p s.nodes }
		newParentNode := s.tree.GetNode(op.NewP.PrntID) 
		if newParentNode.PrntID == s.tree.Tombstone() { // If parent is deleted 
			return
		}
		
		s.tree.Add(op.ChldID, op.NewP)
	case *OpRemove[T]:
		// Checks: requires { n 6= s.root }
		if op.ChldID == s.tree.Root() {
			return
		}
		s.tree.Remove(op.ChldID)
	}
}




// Checks if the state is equal to another state
func (s *State[MD, T]) Equals(other *State[MD, T]) bool {
	treeEq := s.tree.Equals(&other.tree)
	if !treeEq {
		return false
	}
	if s.moveOps.Len() != other.moveOps.Len() {
		return false
	}

	return true
}

func (s *State[MD, T]) DoMoveOp(opMov OpMove[MD, T]) {
	node := s.tree.GetNode(opMov.ChldID) // Get the node
	// Checks: requires { F.mem c s.nodes }
	if node != nil {
		if node.PrntID == s.tree.Tombstone() { // Cannot move deleted node
			return
		}
	} else {
		return
	}
	newPNode := s.tree.GetNode(opMov.NewP.PrntID) 
	// Checks: requires { F.mem np s.nodes }
	if node != nil {
		if newPNode.PrntID == s.tree.Tombstone() { // Cannot move to deleted node
			return
		}
	} else {
		return
	}
	// Checks: requires { c 6= s.root }
	if opMov.ChldID == s.tree.Root() {
		return
	}
	// Checks: requires { c 6= np }
	if opMov.ChldID == opMov.NewP.PrntID {
		return
	}
	// Checks: requires { not (reachability s.parent np c) }
	if isAnc, _ := s.tree.IsAncestor(opMov.NewP.PrntID, opMov.ChldID); isAnc {
		return
	}

	s.tree.Move(opMov.ChldID, opMov.NewP)
	s.moveOps.PushBack(opMov)
}

// Attempted implementation of algorithm from the paper
func (s *State[MD, T]) ApplyMoveOp(opMov OpMove[MD, T]) {
	if s.moveOps.Len() == 0 {
		s.tree.Move(opMov.ChldID, opMov.NewP)
		s.moveOps.PushBack(opMov)
		return
	}

	// Get concurrent operations
	e := s.moveOps.Back()
	for e != nil && e.Value.(OpMove[MD, T]).Timestamp().Compare(opMov.Timestamp()) == 2 {
		e = e.Prev()
	}
	// Get e to the first concurrent operation
	if e == nil {
		e = s.moveOps.Front()
	} else {
		e = e.Next()
	}

	// If no concurrent
	if e == nil {
		s.tree.Move(opMov.ChldID, opMov.NewP)
		s.moveOps.PushBack(opMov)
		return
	}

	// If there are concurrent operations
	
		// If child is the same
	// if conOpMov.ChldID == opMov.ChldID {
	// 	if opMov.ComparePriority(conOpMov) == 1 {
	// 		s.tree.Move(opMov.ChldID, opMov.NewP)
	// 	} 
	// } 
	
	
}


// Up Move Operation
//
// CANT DO IF:
// Exists a concurrent UP-move THAT:
//				- moves the same node, WITH a higher priority

// Down Move Operation
//
// CANT DO IF:
// Exists a concurrent UP-move THAT: 
// 				- that moves the same node 
// 				- OR has a crit-anc-overlap
//
// 				OR 
//
// Exists a concurrent DOWN-move THAT:
// 				- moves the same node 
// 				- OR has a crit-anc-overlap
// 				AND
// 				- has a higher priority