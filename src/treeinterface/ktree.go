package treeinterface

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	k "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/kleppmann"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack" // msgpack is faster and smaller than JSON
)

type KTree[MD any] struct {
	crdt         *k.TreeReplica[MD, *clocks.Lamport]
	crdtMu       sync.RWMutex
	connProv     connection.ConnectionProvider
	totalApplied uint64
}

func (kt *KTree[MD]) GetTotalApplied() uint64 {
	return atomic.LoadUint64(&kt.totalApplied)
}

func NewKTree[MD any](connProv connection.ConnectionProvider) *KTree[MD] {
	kt := &KTree[MD]{crdt: k.NewTreeReplica[MD](nil), connProv: connProv}

	go connProv.HandleBroadcast()
	go connProv.Listen()
	go kt.applyOps(connProv.IncomingOpsChannel())
	return kt
}

// Takes operations from the incoming channel and delivers
// them to the CRDT instantly
func (kt *KTree[MD]) applyOps(ops chan []byte) {
	for {
		opBytes := <-ops

		var op k.OpMove[MD, *clocks.Lamport]
		msgpack.Unmarshal(opBytes, &op)

		kt.crdtMu.Lock()
		kt.crdt.Effect(&op)
		atomic.AddUint64(&kt.totalApplied, 1)
		kt.crdtMu.Unlock()
	}
}

func (kt *KTree[MD]) ConnectionProvider() connection.ConnectionProvider {
	return kt.connProv
}

func (kt *KTree[MD]) Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error) {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	if kt.crdt.GetNode(parentID) == nil {
		return uuid.Nil, errors.New("parent node does not exist")
	}

	id := uuid.New()
	op := kt.crdt.Prepare(id, parentID, metadata)
	if op == nil {
		return uuid.Nil, errors.New("could not prepare operation")
	}

	opBytes, err := msgpack.Marshal(*op)
	if err != nil {
		return uuid.Nil, err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)
	atomic.AddUint64(&kt.totalApplied, 1)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return id, nil
}

func (kt *KTree[MD]) Delete(id uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.Prepare(id, kt.crdt.TombstoneID(), kt.crdt.GetNode(id).Metadata())
	if op == nil {
		return errors.New("could not prepare operation")
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

func (kt *KTree[MD]) Move(id uuid.UUID, newParentID uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}
	if kt.crdt.GetNode(newParentID) == nil {
		return errors.New("new parent node does not exist")
	}

	op := kt.crdt.Prepare(id, newParentID, kt.crdt.GetNode(id).Metadata())
	if op == nil {
		return errors.New("could not prepare operation")
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

func (kt *KTree[MD]) Edit(id uuid.UUID, newMetadata MD) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.Prepare(id, kt.crdt.GetNode(id).ParentID(), newMetadata)
	if op == nil {
		return errors.New("could not prepare operation")
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

func (kt *KTree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	children, bool := kt.crdt.GetChildren(id)
	if !bool {
		return nil, errors.New("node does not exist")
	}
	return children, nil
}

func (kt *KTree[MD]) GetParent(id uuid.UUID) (uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return uuid.Nil, errors.New("node does not exist")
	}
	return node.ParentID(), nil
}

func (kt *KTree[MD]) Root() uuid.UUID {
	return kt.crdt.RootID() // RootID is a constant, so no lock
}

func (kt *KTree[MD]) GetMetadata(id uuid.UUID) (MD, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	var metadata MD
	if node == nil {
		return metadata, errors.New("node does not exist")
	}
	return node.Metadata(), nil
}

func (kt *KTree[MD]) Get(id uuid.UUID) (*tcrdt.TreeNode[MD], error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return nil, errors.New("node does not exist")
	}
	return node, nil
}

func (kt *KTree[MD]) Equals(other *KTree[MD]) bool {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()

	return kt.crdt.Equals(other.crdt)
}
