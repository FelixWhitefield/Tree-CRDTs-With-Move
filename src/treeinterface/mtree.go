package treeinterface

import (
	"errors"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/maram"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"

	//"github.com/shamaton/msgpack/v2" // msgpack is faster and smaller than JSON

	"github.com/vmihailenco/msgpack"

	//msgpack "encoding/json"

	//"google.golang.org/protobuf/proto"
	// This certain encoder is one of the fastest msgpack encoders and decoders
	"container/list"
	"sync"
)

type MTree[MD any] struct {
	crdt     *maram.TreeReplica[MD, *clocks.VectorTimestamp]
	crdtMu   sync.RWMutex
	connProv connection.ConnectionProvider
	opBuffer *list.List
}

func NewMTree[MD any](connProv connection.ConnectionProvider) *MTree[MD] {
	kt := &MTree[MD]{crdt: maram.NewTreeReplica[MD](), connProv: connProv, opBuffer: list.New()}

	go connProv.HandleBroadcast()
	go connProv.Listen()
	go kt.applyOps(connProv.IncomingOpsChannel())

	// Register the custom types for msgpack (They may have already been registered, so we defer the panic)
	kt.RegisterOpMove()
	kt.RegisterOpAdd()
	kt.RegisterOpRemove()

	return kt
}

func (kt *MTree[MD]) RegisterOpMove() {
	defer func() { recover() }()
	msgpack.RegisterExt(1, (*maram.OpMove[MD, *clocks.VectorTimestamp])(nil))
}

func (kt *MTree[MD]) RegisterOpAdd() {
	defer func() { recover() }()
	msgpack.RegisterExt(2, (*maram.OpAdd[MD, *clocks.VectorTimestamp])(nil))
}

func (kt *MTree[MD]) RegisterOpRemove() {
	defer func() { recover() }()
	msgpack.RegisterExt(3, (*maram.OpRemove[*clocks.VectorTimestamp])(nil))
}

func (kt *MTree[MD]) applyOps(ops chan []byte) {
	for {
		opBytes := <-ops

		var op maram.Operation[*clocks.VectorTimestamp]
		msgpack.Unmarshal(opBytes, &op)

		if kt.opBuffer.Len() == 0 {
			kt.opBuffer.PushFront(op)
		} else {
			pos := kt.opBuffer.Front()
			for pos != nil && op.Timestamp().Compare(pos.Value.(maram.Operation[*clocks.VectorTimestamp]).Timestamp()) == 1 {
				pos = pos.Next()
			}
			if pos == nil {
				if op.Timestamp().Same(kt.opBuffer.Back().Value.(maram.Operation[*clocks.VectorTimestamp]).Timestamp()) {
				} else {
					kt.opBuffer.PushBack(op)
				}
			} else {
				if op.Timestamp().Same(pos.Value.(maram.Operation[*clocks.VectorTimestamp]).Timestamp()) {
				} else {
					kt.opBuffer.InsertBefore(op, pos)
				}
			}
		}

		front := kt.opBuffer.Front()
		kt.crdtMu.Lock()
		for front != nil && (front.Value.(maram.Operation[ *clocks.VectorTimestamp]).Timestamp().CausallyReady(kt.crdt.CurrentTime()) || 
			front.Value.(maram.Operation[*clocks.VectorTimestamp]).Timestamp().Compare(kt.crdt.CurrentTime()) == -1) {
			if front.Value.(maram.Operation[ *clocks.VectorTimestamp]).Timestamp().CausallyReady(kt.crdt.CurrentTime()) {
				opToApp := front.Value.(maram.Operation[*clocks.VectorTimestamp])
				kt.crdt.Effect(opToApp)
			}
			kt.opBuffer.Remove(front)
			front = kt.opBuffer.Front()
		}
		kt.crdtMu.Unlock()
	}
}

func (kt *MTree[MD]) ConnectionProvider() connection.ConnectionProvider {
	return kt.connProv
}

func (kt *MTree[MD]) Equals(other *MTree[MD]) bool {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()

	return kt.crdt.State().Equals(other.crdt.State())
}

func (kt *MTree[MD]) Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error) {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	if kt.crdt.GetNode(parentID) == nil {
		return uuid.Nil, errors.New("parent node does not exist")
	}

	id := uuid.New()
	op := kt.crdt.PrepareAdd(parentID, metadata)

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return uuid.Nil, err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return id, nil
}

func (kt *MTree[MD]) Delete(id uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.PrepareRemove(id)

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op) // Apply the operation to the state (After it is successfully encoded)

	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *MTree[MD]) Move(id uuid.UUID, newParentID uuid.UUID) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}
	if kt.crdt.GetNode(newParentID) == nil {
		return errors.New("new parent node does not exist")
	}

	op := kt.crdt.PrepareMove(id, newParentID, kt.crdt.GetNode(id).Metadata())

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *MTree[MD]) Edit(id uuid.UUID, newMetadata MD) error {
	kt.crdtMu.Lock()
	defer kt.crdtMu.Unlock()

	node := kt.crdt.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.crdt.PrepareMove(id, kt.crdt.GetNode(id).ParentID(), newMetadata)

	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.crdt.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *MTree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	children, bool := kt.crdt.GetChildren(id)
	if !bool {
		return nil, errors.New("node does not exist")
	}
	return children, nil
}

func (kt *MTree[MD]) GetParent(id uuid.UUID) (uuid.UUID, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return uuid.Nil, errors.New("node does not exist")
	}
	return node.ParentID(), nil
}

func (kt *MTree[MD]) Root() uuid.UUID {
	return kt.crdt.RootID() // RootID is a constant, so no lock
}

func (kt *MTree[MD]) GetMetadata(id uuid.UUID) (MD, error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	var metadata MD
	if node == nil {
		return metadata, errors.New("node does not exist")
	}
	return node.Metadata(), nil
}

func (kt *MTree[MD]) Get(id uuid.UUID) (*tcrdt.TreeNode[MD], error) {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()
	node := kt.crdt.GetNode(id)
	if node == nil {
		return nil, errors.New("node does not exist")
	}
	return node, nil
}
