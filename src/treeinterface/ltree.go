package treeinterface

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"

	//"github.com/theodesp/go-heaps/fibonacci"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/lumina"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack" // msgpack is faster and smaller than JSON
)

type LTree[MD any] struct {
	crdt         *lumina.TreeReplica[MD, *clocks.VectorTimestamp]
	crdtMu       sync.RWMutex
	connProv     connection.ConnectionProvider
	totalApplied uint64
}

func NewLTree[MD any](connProv connection.ConnectionProvider, optimisedBuffer bool) *LTree[MD] {
	kt := &LTree[MD]{crdt: lumina.NewTreeReplica[MD](), connProv: connProv}

	go connProv.HandleBroadcast()
	go connProv.Listen()
	if optimisedBuffer {
		go kt.applyOpsSkip(connProv.IncomingOpsChannel())
	} else {
		go kt.applyOps(connProv.IncomingOpsChannel())
	}

	// Register the custom types for msgpack (They may have already been registered, so we defer the panic)
	kt.RegisterOpMove()
	kt.RegisterOpAdd()
	kt.RegisterOpRemove()

	return kt
}

func (kt *LTree[MD]) GetTotalApplied() uint64 {
	return atomic.LoadUint64(&kt.totalApplied)
}

func (kt *LTree[MD]) RegisterOpMove() {
	defer func() { recover() }()
	msgpack.RegisterExt(1, (*lumina.OpMove[MD, *clocks.VectorTimestamp])(nil))
}

func (kt *LTree[MD]) RegisterOpAdd() {
	defer func() { recover() }()
	msgpack.RegisterExt(2, (*lumina.OpAdd[MD, *clocks.VectorTimestamp])(nil))
}

func (kt *LTree[MD]) RegisterOpRemove() {
	defer func() { recover() }()
	msgpack.RegisterExt(3, (*lumina.OpRemove[*clocks.VectorTimestamp])(nil))
}

func (kt *LTree[MD]) applyOpsSkip(ops chan []byte) {
	skip := 100
	opBuffer := make([]lumina.Operation[*clocks.VectorTimestamp], 0)
	for {
		opBytes := <-ops

		var op lumina.Operation[*clocks.VectorTimestamp]
		msgpack.Unmarshal(opBytes, &op)

		i := 0
		if len(opBuffer) > skip {
			for i < len(opBuffer) {
				if compare := op.Timestamp().Compare(opBuffer[i].Timestamp()); compare == -1 || compare == 2 {
					break
				}
				i += skip
			}
			i -= skip
		}
		if i < 0 {
			i = 0
		}
		for i < len(opBuffer) {
			if compare := op.Timestamp().Compare(opBuffer[i].Timestamp()); compare == -1 {
				break
			}
			i++
		}
		opBuffer = append(opBuffer, nil)
		copy(opBuffer[i+1:], opBuffer[i:])
		opBuffer[i] = op

		min := opBuffer[0]
		kt.crdtMu.Lock()
		for min != nil {
			causallyReady := min.Timestamp().CausallyReady(kt.crdt.CurrentTime())
			compare := min.Timestamp().Compare(kt.crdt.CurrentTime())
			if causallyReady {
				kt.crdt.Effect(min)
				atomic.AddUint64(&kt.totalApplied, 1)
				opBuffer = opBuffer[1:]
			} else if compare == -1 || compare == 0 {
				opBuffer = opBuffer[1:]
			} else {
				break
			}
			if len(opBuffer) == 0 {
				break
			}
			min = opBuffer[0]
		}
		kt.crdtMu.Unlock()
	}
}

// Takes operations from the incoming channel and delivers
// them in a causal order to the CRDT
func (kt *LTree[MD]) applyOps(ops chan []byte) {
	opBuffer := list.New()
	for {
		opBytes := <-ops

		var op lumina.Operation[*clocks.VectorTimestamp]
		msgpack.Unmarshal(opBytes, &op)

		if opBuffer.Len() == 0 {
			opBuffer.PushFront(op)
		} else {
			pos := opBuffer.Front()
			for pos != nil {
				compare := op.Timestamp().Compare(pos.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp())
				if compare == -1 {
					break
				}
				pos = pos.Next()
			}
			if pos == nil {
				opBuffer.PushBack(op)
			} else if compare := op.Timestamp().Compare(pos.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp()); compare == 0 {
				continue
			} else {
				opBuffer.InsertBefore(op, pos)
			}
		}

		item := opBuffer.Front()

		kt.crdtMu.Lock()
		for item != nil {
			op := item.Value.(lumina.Operation[*clocks.VectorTimestamp])
			causallyReady := op.Timestamp().CausallyReady(kt.crdt.CurrentTime())
			if causallyReady {
				opToApp := item.Value.(lumina.Operation[*clocks.VectorTimestamp])
				kt.crdt.Effect(opToApp)
				atomic.AddUint64(&kt.totalApplied, 1)
				opBuffer.Remove(item)
			} else if compare := op.Timestamp().Compare(kt.crdt.CurrentTime()); compare == -1 || compare == 0 {
				opBuffer.Remove(item)
			} else {
				break
			}
			item = opBuffer.Front()
		}
		kt.crdtMu.Unlock()
	}
}

func (kt *LTree[MD]) ConnectionProvider() connection.ConnectionProvider {
	return kt.connProv
}

func (kt *LTree[MD]) Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error) {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	if kt.crdt.GetNode(parentID) == nil {
		return uuid.Nil, errors.New("parent node does not exist")
	}

	id := uuid.New()
	op := kt.crdt.PrepareAdd(id, parentID, metadata)
	if op == nil {
		return uuid.Nil, errors.New("error preparing add")
	}

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return uuid.Nil, err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)
	atomic.AddUint64(&kt.totalApplied, 1)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return id, nil
}

func (kt *LTree[MD]) Delete(id uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.PrepareRemove(id)
	if op == nil {
		return errors.New("error preparing remove")
	}

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)
	atomic.AddUint64(&kt.totalApplied, 1)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *LTree[MD]) Move(id uuid.UUID, newParentID uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}
	if kt.crdt.GetNode(newParentID) == nil {
		return errors.New("new parent node does not exist")
	}

	op := kt.crdt.PrepareMove(id, newParentID, node.Metadata())
	if op == nil {
		return errors.New("error preparing move")
	}

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)
	atomic.AddUint64(&kt.totalApplied, 1)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *LTree[MD]) Edit(id uuid.UUID, newMetadata MD) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.PrepareMove(id, kt.crdt.GetNode(id).ParentID(), newMetadata)
	if op == nil {
		return errors.New("error preparing edit")
	}

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)
	atomic.AddUint64(&kt.totalApplied, 1)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *LTree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	children, bool := kt.crdt.GetChildren(id)
	if !bool {
		return nil, errors.New("node does not exist")
	}
	return children, nil
}

func (kt *LTree[MD]) GetParent(id uuid.UUID) (uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return uuid.Nil, errors.New("node does not exist")
	}
	return node.ParentID(), nil
}

func (kt *LTree[MD]) Root() uuid.UUID {
	return kt.crdt.RootID() // RootID is a constant, so no lock
}

func (kt *LTree[MD]) GetMetadata(id uuid.UUID) (MD, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	var metadata MD
	if node == nil {
		return metadata, errors.New("node does not exist")
	}
	return node.Metadata(), nil
}

func (kt *LTree[MD]) Get(id uuid.UUID) (*tcrdt.TreeNode[MD], error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return nil, errors.New("node does not exist")
	}
	return node, nil
}

func (kt *LTree[MD]) Equals(other *LTree[MD]) bool {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()

	return kt.crdt.State().Equals(other.crdt.State())
}