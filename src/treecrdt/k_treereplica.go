package treecrdt 

// `TreeReplica` is a replica of a tree CRDT. It contains the state of the replica and the clock of the replica.
//
// The replica is an implementation of a op-based CRDT, and contains `prepare` and `effect` methods
// This is a layer above the `State` struct, which contains the actual CRDT state
// The replica is responsible for applying operations to the state, and for generating operations

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/google/uuid"
)

type TreeReplica[MD Metadata, T opTimestamp[T]] struct {
	state State[MD, T] // contains the state of the replica
	clock c.Clock[T] // contains current time of replica (including actorID)
}

// Returns a new TreeReplica with a random actorID, using the Lamport clock
func NewTreeReplica[MD Metadata]() TreeReplica[MD, *c.Lamport] {
	return TreeReplica[MD, *c.Lamport]{state: NewState[MD, *c.Lamport](), clock: c.NewLamport()}
}

// Returns a new TreeReplica with the given actorID, using the Lamport clock
func NewTreeReplicaWithID[MD Metadata](id uuid.UUID) TreeReplica[MD, *c.Lamport] {
	return TreeReplica[MD, *c.Lamport]{state: NewState[MD, *c.Lamport](), clock: c.NewLamport(id)}
}

func (tr *TreeReplica[MD, T]) ActorID() uuid.UUID {
	return tr.clock.ActorID()
}

// The `prepare` method for the op-based CRDt, prepares an operation for the replica.
func (tr *TreeReplica[MD, T]) Prepare(id uuid.UUID, newP uuid.UUID, metadata MD) *OpMove[MD, T] {
	tclock := tr.clock.Tick()
	return NewOpMove(tclock, newP, id, metadata)
}

// The `effect` method for the op-based CRDt, applies an operation to the replica.
// This creates the effect of the operation on the replica.
func (tr *TreeReplica[MD, T]) Effect(op OpMove[MD, T]) {
	tr.state.ApplyOp(op)
	tr.clock.Merge(op.Timestamp().(T))
}
