package lumina

// Represents the state of the CRDT for a single replica

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

type TreeReplica[MD any, T opTimestamp[T]] struct {
	state    State[MD, T]
	clock    c.Clock[T]
	priotity *c.Lamport
}

func NewTreeReplica[MD any](ids ...uuid.UUID) *TreeReplica[MD, *c.VectorTimestamp] {
	var id uuid.UUID
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = uuid.New()
	}
	return &TreeReplica[MD, *c.VectorTimestamp]{state: *NewLState[MD, *c.VectorTimestamp](), clock: c.NewVectorClock(id), priotity: c.NewLamport(id)}
}

func (tr *TreeReplica[MD, T]) RootID() uuid.UUID {
	return tr.state.tree.Root()
}

func (tr *TreeReplica[MD, T]) ActorID() uuid.UUID {
	return tr.clock.ActorID()
}

func (tr *TreeReplica[MD, T]) CurrentTime() T {
	return tr.clock.Timestamp()
}

func (tr *TreeReplica[MD, T]) GetChildren(u uuid.UUID) ([]uuid.UUID, bool) {
	return tr.state.tree.GetChildren(u)
}

func (tr *TreeReplica[MD, T]) GetNode(u uuid.UUID) *treecrdt.TreeNode[MD] {
	return tr.state.tree.GetNode(u)
}

// Prepares an add operation
func (tr *TreeReplica[MD, T]) PrepareAdd(childId uuid.UUID, parentId uuid.UUID, metadata MD) *OpAdd[MD, T] {
	if !tr.state.tree.Contains(parentId) {
		return nil
	}
	return &OpAdd[MD, T]{Timestmp: tr.clock.Tick(), ChldID: childId, NewP: &treecrdt.TreeNode[MD]{PrntID: parentId, Meta: metadata}}
}

// Prepares a remove operation
func (tr *TreeReplica[MD, T]) PrepareRemove(id uuid.UUID) *OpRemove[T] {
	if !tr.state.tree.Contains(id) {
		return nil
	}
	return &OpRemove[T]{Timestmp: tr.clock.Tick(), ChldID: id}
}

// Prepares a move operation
func (tr *TreeReplica[MD, T]) PrepareMove(id uuid.UUID, newP uuid.UUID, metadata MD) *OpMove[MD, T] {
	childIsAnc, _ := tr.state.tree.IsAncestor(newP, id)
	if !tr.state.tree.Contains(id) || !tr.state.tree.Contains(newP) || id == newP || id == tr.state.tree.Root() || childIsAnc {
		return nil
	}
	return &OpMove[MD, T]{Timestmp: tr.clock.Tick(), ChldID: id, NewP: &treecrdt.TreeNode[MD]{PrntID: newP, Meta: metadata}, Priotity: *tr.priotity.Tick()}
}

// Applies an operation to the tree, and updates the clock
func (tr *TreeReplica[MD, T]) Effect(op Operation[T]) {
	tr.clock.Merge(op.Timestamp())

	switch op := op.(type) { // If op is of type OpAdd, update priority
	case *OpMove[MD, T]:
		tr.priotity.Merge(&op.Priotity)
	}

	tr.state.ApplyOp(op)
}

// Applies a list of operations to the tree
func (tr *TreeReplica[MD, T]) Effects(op []Operation[T]) {
	for _, o := range op {
		tr.Effect(o)
	}
}

func (tr *TreeReplica[MD, T]) State() *State[MD, T] {
	return &tr.state
}

func (tr *TreeReplica[MD, T]) Equals(other *TreeReplica[MD, T]) bool {
	return tr.state.Equals(&other.state)
}
