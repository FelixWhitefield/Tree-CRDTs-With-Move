package kleppmann

// `TreeReplica` is a replica of the tree CRDT. It contains the state of the replica and the clock of the replica.
//
// This struct is the CRDT for a certain actor, and contains the state of the CRDT for that actor.
//
// The replica is an implementation of a op-based CRDT, and contains `prepare` and `effect` methods
// This is a layer above the `State` struct, which contains the actual CRDT state
// The replica is responsible for applying operations to the state, and for generating operations

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

type TreeReplica[MD any, T opTimestamp[T]] struct {
	state                     State[MD, T]    // contains the state of the replica
	clock                     c.Clock[T]      // contains current time of replica (including actorID)
	latest_timestamp_by_actor map[uuid.UUID]T // contains the latest timestamp of each actor
}

// Returns a new TreeReplica with a random actorID, using the Lamport clock
// func NewTreeReplica[MD Metadata]() *TreeReplica[MD, *c.Lamport] {
// 	return &TreeReplica[MD, *c.Lamport]{state: NewState[MD, *c.Lamport](), clock: c.NewLamport(), latest_timestamp_by_actor: make(map[uuid.UUID]*c.Lamport)}
// }

// Returns a new TreeReplica with the given actorID, using the Lamport clock
//
//	func NewTreeReplicaWithID[MD Metadata](id uuid.UUID) *TreeReplica[MD, *c.Lamport] {
//		return &TreeReplica[MD, *c.Lamport]{state: NewState[MD, *c.Lamport](), clock: c.NewLamport(id), latest_timestamp_by_actor: make(map[uuid.UUID]*c.Lamport)}
//	}
func NewTreeReplica[MD any](conf *TNConflict[MD], ids ...uuid.UUID) *TreeReplica[MD, *c.Lamport] {
	var id uuid.UUID
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = uuid.New()
	}
	return &TreeReplica[MD, *c.Lamport]{state: *NewState[MD, *c.Lamport](conf), clock: c.NewLamport(id), latest_timestamp_by_actor: make(map[uuid.UUID]*c.Lamport)}
}

func (tr *TreeReplica[MD, T]) RootID() uuid.UUID {
	return tr.state.tree.Root()
}

func (tr *TreeReplica[MD, T]) TombstoneID() uuid.UUID {
	return tr.state.tree.Tombstone()
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

func (tr *TreeReplica[MD, T]) GetNode(u uuid.UUID) *TreeNode[MD] {
	return tr.state.tree.GetNode(u)
}

// The `prepare` method for the op-based CRDTs, prepares an operation for the replica.
func (tr *TreeReplica[MD, T]) Prepare(id uuid.UUID, newP uuid.UUID, metadata MD) *OpMove[MD, T] {
	childIsAnc, _ := tr.state.tree.IsAncestor(newP, id)
	if id == tr.state.tree.Root() || !tr.state.tree.Contains(newP) || childIsAnc {
		return nil
	}
	return NewOpMove(tr.clock.Tick(), newP, id, metadata)
}

// The `effect` method for the op-based CRDTs, applies an operation to the replica.
// This creates the effect of the operation on the replica.
func (tr *TreeReplica[MD, T]) Effect(op *OpMove[MD, T]) {
	tr.clock.Merge(op.Timestamp())

	id := op.Timestamp().ActorID()
	// if the latest timestamp of the actor is less than the timestamp of the operation
	if latest, exist := tr.latest_timestamp_by_actor[id]; !exist || latest.Compare(op.Timestamp()) == -1 {
		tr.latest_timestamp_by_actor[id] = op.Timestamp().Clone()
	}

	tr.state.ApplyOp(op)
}

// Applies multiple operations to the replica
func (tr *TreeReplica[MD, T]) Effects(ops []*OpMove[MD, T]) {
	for _, op := range ops {
		tr.Effect(op)
	}
}

func (tr *TreeReplica[MD, T]) CausallyStableThreshold() *T {
	var min *T
	for _, timestamp := range tr.latest_timestamp_by_actor {
		if min == nil || timestamp.Compare(*min) == -1 {
			min = &timestamp
		}
	}
	return min
}

func (tr *TreeReplica[MD, T]) TruncateLog() {
	threshold := tr.CausallyStableThreshold()
	if threshold != nil {
		tr.state.TruncateLogBefore(*threshold)
	}
}

// Should be used to truncate log from outside
// The log should only be truncated if state has seen all other peers
// func (tr *TreeReplica[MD, T]) TruncateLogMinPeers(n int) {
// 	if len(tr.latest_timestamp_by_actor) >= n {
// 		tr.TruncateLog()
// 	}
// }
