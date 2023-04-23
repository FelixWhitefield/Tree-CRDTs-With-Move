package treeinterface

import (
	"errors"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/k"
	"github.com/google/uuid"
	"github.com/shamaton/msgpack/v2" // msgpack is faster and smaller than JSON
	//"google.golang.org/protobuf/proto"
	// This certain encoder is one of the fastest msgpack encoders and decoders
)

type KTree[MD any] struct {
	tree     *k.TreeReplica[MD, *clocks.Lamport]
	connProv connection.ConnectionProvider
}

func NewKTree[MD any](connProv connection.ConnectionProvider) *KTree[MD] {
	go connProv.HandleBroadcast()
	go connProv.Listen()

	kt := &KTree[MD]{tree: k.NewTreeReplica[MD](nil), connProv: connProv}
	go kt.ApplyOps(connProv.IncomingOpsChannel())
	return kt
}

func (kt *KTree[MD]) ApplyOps(ops chan []byte) {
	for {
		opBytes := <-ops

		var op k.OpMove[MD, *clocks.Lamport]
		msgpack.Unmarshal(opBytes, &op)
		kt.tree.Effect(&op)
	}
}

func (kt *KTree[MD]) ConnectionProvider() connection.ConnectionProvider {
	return kt.connProv
}

func (kt *KTree[MD]) Insert(parentID uuid.UUID, metadata MD) (uuid.UUID, error) {
	if kt.tree.GetNode(parentID) == nil {
		return uuid.Nil, errors.New("parent node does not exist")
	}

	id := uuid.New()
	op := kt.tree.Prepare(id, parentID, metadata)
	opBytes, err := msgpack.Marshal(*op)
	if err != nil {
		return uuid.Nil, err
	}

	kt.tree.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return id, nil
}

func (kt *KTree[MD]) Delete(id uuid.UUID) error {
	node := kt.tree.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.tree.Prepare(id, kt.tree.TombstoneID(), kt.tree.GetNode(id).Metadata())
	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.tree.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *KTree[MD]) Move(id uuid.UUID, newParentID uuid.UUID) error {
	node := kt.tree.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}
	if kt.tree.GetNode(newParentID) == nil {
		return errors.New("new parent node does not exist")
	}

	op := kt.tree.Prepare(id, newParentID, kt.tree.GetNode(id).Metadata())
	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.tree.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *KTree[MD]) Edit(id uuid.UUID, newMetadata MD) error {
	node := kt.tree.GetNode(id)
	if node == nil {
		return errors.New("node does not exist")
	}

	op := kt.tree.Prepare(id, kt.tree.GetNode(id).ParentID(), newMetadata)
	opBytes, err := msgpack.Marshal(op)
	if err != nil {
		return err
	}

	kt.tree.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
	kt.connProv.BroadcastChannel() <- opBytes // Broadcast Op

	return nil
}

func (kt *KTree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, error) {
	children, bool := kt.tree.GetChildren(id)
	if !bool {
		return nil, errors.New("node does not exist")
	}
	return children, nil
}

func (kt *KTree[MD]) GetParent(id uuid.UUID) (uuid.UUID, error) {
	node := kt.tree.GetNode(id)
	if node == nil {
		return uuid.Nil, errors.New("node does not exist")
	}
	return node.ParentID(), nil
}

func (kt *KTree[MD]) Root() uuid.UUID {
	return kt.tree.RootID()
}

func (kt *KTree[MD]) Get(id uuid.UUID) (MD, error) {
	node := kt.tree.GetNode(id)
	var metadata MD
	if node == nil {
		return metadata, errors.New("node does not exist")
	}
	return node.Metadata(), nil
}
