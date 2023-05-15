package lumina

// This tree implementation was created to allow Lumina to correctly function
//
// The generic tree implementation does not work for Lumina, due to Lumina
// not logging all operations. Therefore, delete operations can cause issues
// 
// This tree implementation is based on the generic tree implementation, but
// instead of a tombstone node, this tree implementation uses a tombstone flag
// on the node itself. This allows the tree to correctly function with Lumina

import (
	"bytes"
	"errors"
	"fmt"

	//"fmt"
	//"log"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

var (
	RootUUID      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

type LuminaNode[MD any] struct {
	Node *treecrdt.TreeNode[MD]
	Tombstone bool
}

type Tree[MD any] struct {
	nodes    map[uuid.UUID]*LuminaNode[MD]      // node id -> tree node
	children map[uuid.UUID]map[uuid.UUID]bool // node id -> set child id
}

func NewTree[MD any]() *Tree[MD] {
	tree := Tree[MD]{nodes: make(map[uuid.UUID]*LuminaNode[MD]), children: make(map[uuid.UUID]map[uuid.UUID]bool)}
	tree.nodes[RootUUID] = &LuminaNode[MD]{Node: &treecrdt.TreeNode[MD]{PrntID: RootUUID}}
	return &tree
}

func (t *Tree[MD]) Root() uuid.UUID {
	return RootUUID
}

func (t *Tree[MD]) Add(id uuid.UUID, node *treecrdt.TreeNode[MD]) error {
	if _, exists := t.nodes[id]; exists { // if already exists, there may be duplicate data
		return errors.New("node already exists")
	}
	if _, exists := t.nodes[node.PrntID]; !exists {
		return errors.New("parent does not exist")
	}
	t.nodes[id] = &LuminaNode[MD]{Node: node, Tombstone: false}
	if _, exists := t.children[node.PrntID]; !exists { // If no children set, create one
		t.children[node.PrntID] = make(map[uuid.UUID]bool)
		t.children[node.PrntID][id] = true
	} else { // otherwise add to existing
		t.children[node.PrntID][id] = true
	}
	return nil
}

func (t *Tree[MD]) Remove(id uuid.UUID) error {
	if _, exists := t.nodes[id]; !exists {
		return errors.New("node does not exist")
	}
	if id == RootUUID {
		return errors.New("cannot remove root node")
	}
	t.nodes[id].Tombstone = true
	return nil
}

func (t *Tree[MD]) Move(id uuid.UUID, node *treecrdt.TreeNode[MD]) error {
	if _, exists := t.nodes[id]; !exists {
		return errors.New("node does not exist")
	}	
	if _, exists := t.nodes[node.PrntID]; !exists {
		return errors.New("parent does not exist")
	}
	if id == RootUUID {
		return errors.New("cannot move root node")
	}
	parentID := t.nodes[id].Node.PrntID
	delete(t.children[parentID], id) // Remove from old parent
	isTombstone := t.nodes[id].Tombstone
	delete(t.nodes, id)

	t.nodes[id] = &LuminaNode[MD]{Node: node, Tombstone: isTombstone}
	if _, exists := t.children[node.PrntID]; !exists { // If no children set, create one
		t.children[node.PrntID] = make(map[uuid.UUID]bool)
		t.children[node.PrntID][id] = true
	} else { // otherwise add to existing
		t.children[node.PrntID][id] = true
	}
	return nil
}


func (t *Tree[MD]) GetNode(id uuid.UUID) *treecrdt.TreeNode[MD] {
	if !t.InTree(id) {
		return nil
	}
	node := t.nodes[id]
	return node.Node
}

func (t *Tree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, bool) {
	if !t.InTree(id) {
		return nil, false
	}
	children, exists := t.children[id]
	result := make([]uuid.UUID, 0, len(children))
	for k := range children {
		if t.nodes[k].Tombstone {
			continue
		}
		result = append(result, k)
	}
	return result, exists
}

// If the node is in the tree (may be tombstoned)
func (t *Tree[MD]) WithinTree(id uuid.UUID) *LuminaNode[MD] {
	node := t.nodes[id]
	return node
}

// If the node is in the tree and not tombstoned
func (t *Tree[MD]) InTree(id uuid.UUID) bool {
	node, exists := t.nodes[id]
	if !exists {
		return false
	}
	for id != RootUUID {
		if node.Tombstone {
			return false
		}
		id = node.Node.PrntID
		node = t.nodes[id]
	}
	return true
}

func (t *Tree[MD]) Contains(id uuid.UUID) bool {
	return t.InTree(id)
}

// Checks if node is ancestor of other node.
// E.g. If ancID can be reached by following the parent pointers from childID.
func (t *Tree[MD]) IsAncestor(childID uuid.UUID, ancID uuid.UUID) (bool, error) {
	if _, exists := t.nodes[childID]; !exists {
		return false, errors.New("child node does not exist")
	}
	if _, exists := t.nodes[ancID]; !exists {
		return false, errors.New("ancestor node does not exist")
	}
	for childID != RootUUID {
		childID = t.nodes[childID].Node.PrntID
		if bytes.Equal(childID[:], ancID[:]) {
			return true, nil
		}
	}
	return false, nil
}

func (t *Tree[MD]) Equals(other *Tree[MD]) bool {
	if len(t.nodes) != len(other.nodes) {
		return false
	}
	for k, v := range t.nodes {
		if !v.Node.Equals(other.nodes[k].Node) {
			return false
		}
	}
	return true
}

func (t *Tree[MD]) String() string {
	var buf bytes.Buffer
	for k, v := range t.nodes {
		buf.WriteString(fmt.Sprintf("%v: %v\n", k, v.Node))
	}
	return buf.String()
}