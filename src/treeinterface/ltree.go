package treeinterface

import (
	"container/list"
	"errors"
	"sync"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	tcrdt "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/lumina"
	rb "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack" // msgpack is faster and smaller than JSON
	"container/heap"
)

type OperationQueue []lumina.Operation[*clocks.VectorTimestamp]

func (pq OperationQueue) Len() int { return len(pq) }

func (pq OperationQueue) Less(i, j int) bool {
	// Compare the timestamps of the operations to determine priority
	return pq[i].Timestamp().Compare(pq[j].Timestamp()) == -1
}

func (pq OperationQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *OperationQueue) Push(x interface{}) {
	*pq = append(*pq, x.(lumina.Operation[*clocks.VectorTimestamp]))
}

func (pq *OperationQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type LTree[MD any] struct {
	crdt     *lumina.TreeReplica[MD, *clocks.VectorTimestamp]
	crdtMu   sync.RWMutex
	connProv connection.ConnectionProvider
	opBuffer *list.List
}

func NewLTree[MD any](connProv connection.ConnectionProvider) *LTree[MD] {
	kt := &LTree[MD]{crdt: lumina.NewTreeReplica[MD](), connProv: connProv, opBuffer: list.New()}

	go connProv.HandleBroadcast()
	go connProv.Listen()
	go kt.applyOps(connProv.IncomingOpsChannel())

	// Register the custom types for msgpack (They may have already been registered, so we defer the panic)
	kt.RegisterOpMove()
	kt.RegisterOpAdd()
	kt.RegisterOpRemove()

	return kt
}

func (kt *LTree[MD]) GetBufLen() int {
	return kt.opBuffer.Len()
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

func (kt *LTree[MD]) applyOpsHeap(ops chan []byte) {
	pq := &OperationQueue{}
	heap.Init(pq)

	for {
		opBytes := <-ops

		var op lumina.Operation[*clocks.VectorTimestamp]
		msgpack.Unmarshal(opBytes, &op)

		heap.Push(pq, op)

		kt.crdtMu.Lock()
		for pq.Len() > 0 {
			causallyReady := (*pq)[0].Timestamp().CausallyReady(kt.crdt.CurrentTime())
			compare := (*pq)[0].Timestamp().Compare(kt.crdt.CurrentTime())
			if causallyReady {
				opToApp := heap.Pop(pq).(lumina.Operation[*clocks.VectorTimestamp])
				kt.crdt.Effect(opToApp)
			} else if compare == -1 || compare == 0 {
				heap.Pop(pq)
			} else {
				break
			}
		}
		kt.crdtMu.Unlock()
	}	
}


// Takes operations from the incoming channel and delivers
// them in a causal order to the CRDT
func (kt *LTree[MD]) applyOps(ops chan []byte) {
	for {
		opBytes := <-ops

		var op lumina.Operation[*clocks.VectorTimestamp]
		msgpack.Unmarshal(opBytes, &op)

		if kt.opBuffer.Len() == 0 {
			kt.opBuffer.PushFront(op)
		} else {
			pos := kt.opBuffer.Front()
			for pos != nil && op.Timestamp().Compare(pos.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp()) == 1 {
				pos = pos.Next()
			}
			if pos == nil {
				kt.opBuffer.PushBack(op)
			} else {
				kt.opBuffer.InsertBefore(op, pos)
			}
		}

		item := kt.opBuffer.Front()

		kt.crdtMu.Lock()
		for item != nil {
			op := item.Value.(lumina.Operation[*clocks.VectorTimestamp])
			causallyReady := op.Timestamp().CausallyReady(kt.crdt.CurrentTime())
			compare := op.Timestamp().Compare(kt.crdt.CurrentTime())
			if causallyReady {
				opToApp := item.Value.(lumina.Operation[*clocks.VectorTimestamp])
				kt.crdt.Effect(opToApp)
				kt.opBuffer.Remove(item)
			} else if compare == -1 || compare == 0 {
				kt.opBuffer.Remove(item)
			} else {
				break
			}
			item = kt.opBuffer.Front()
		}
		kt.crdtMu.Unlock()
		

	

		// for ; item != nil && (item.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp().CausallyReady(kt.crdt.CurrentTime()) ||
		// 	item.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp().Compare(kt.crdt.CurrentTime()) == -1 ||
		// 	item.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp().Compare(kt.crdt.CurrentTime()) == 0); item = kt.opBuffer.Front() {
		// 		if item.Value.(lumina.Operation[*clocks.VectorTimestamp]).Timestamp().CausallyReady(kt.crdt.CurrentTime()) {
		// 			opToApp := item.Value.(lumina.Operation[*clocks.VectorTimestamp])
		// 			kt.crdt.Effect(opToApp)
		// 		}
		// 	kt.opBuffer.Remove(item)

		// }
	}
}


func (kt *LTree[MD]) ConnectionProvider() connection.ConnectionProvider {
	return kt.connProv
}

func (kt *LTree[MD]) Equals(other *LTree[MD]) bool {
	kt.crdtMu.RLock()
	defer kt.crdtMu.RUnlock()

	return kt.crdt.State().Equals(other.crdt.State())
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

	kt.crdt.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
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

	kt.crdt.Effect(op)                        // Apply the operation to the state (After it is successfully encoded)
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
