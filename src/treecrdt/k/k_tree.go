package k

// `Tree` holds the current state of the tree.
//
// The `nodes` map represents the triples indexed by child id.
// The `children` map provides a quick lookup of a node's children.
//
// RootUUID and TombstoneUUID are special nodes that are always present.
// They are set manually as to ensure they are the same across all replicas.

import (
	//"errors"
	"bytes"
	"fmt"

	"github.com/google/uuid"
)

var (
	RootUUID      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	TombstoneUUID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

type Tree[MD any] struct {
	nodes    map[uuid.UUID]*TreeNode[MD] // node id -> tree node
	children map[uuid.UUID][]uuid.UUID   // node id -> []child id
}

// Creates a new tree with the root and tombstone nodes
// This is the proper way to create a new tree.
func NewTree[MD any]() *Tree[MD] {
	tree := Tree[MD]{nodes: make(map[uuid.UUID]*TreeNode[MD]), children: make(map[uuid.UUID][]uuid.UUID)}
	tree.nodes[RootUUID] = &TreeNode[MD]{}
	tree.nodes[TombstoneUUID] = &TreeNode[MD]{}
	return &tree
}

func (t *Tree[MD]) Root() uuid.UUID {
	return RootUUID
}

func (t *Tree[MD]) Tombstone() uuid.UUID {
	return TombstoneUUID
}

// Returns the node with the given id. Returns false if the node does not exist.
func (t *Tree[MD]) GetNode(id uuid.UUID) *TreeNode[MD] {
	node := t.nodes[id]
	return node
}

// Returns the children of the node with the given id. Returns false if the node does not exist.
func (t *Tree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, bool) {
	children, exists := t.children[id]
	return children, exists
}

// Adds a node to the tree.
// Errors if the node already exists.
// Errors if the parent node does not exist.
func (t *Tree[MD]) Add(id uuid.UUID, node *TreeNode[MD]) error {
	if _, exists := t.nodes[id]; exists { // if already exists, there may be duplicate data
		return IDAlreadyExistsError{id: id}
	}
	if _, exists := t.nodes[node.PrntID]; !exists {
		return MissingNodeIDError{id: node.PrntID}
	}

	t.nodes[id] = node
	if _, exists := t.children[node.PrntID]; !exists {
		t.children[node.PrntID] = []uuid.UUID{id}
	} else {
		t.children[node.PrntID] = append(t.children[node.PrntID], id)
	}
	return nil
}

// Removes a node from the tree. Doesn't remove the corresponding children entry.
// Errors if the node does not exist.
// Errors if the node is the root or deleted node.
func (t *Tree[MD]) Remove(id uuid.UUID) error {
	if _, exists := t.nodes[id]; !exists {
		return nil // already removed
	}
	if id == RootUUID || id == TombstoneUUID {
		return InvalidNodeDeletionError{id: id}
	}

	parentID := t.nodes[id].PrntID
	for i, childID := range t.children[parentID] { // remove child from parent's children
		if childID == id {
			t.children[parentID] = append(t.children[parentID][:i], t.children[parentID][i+1:]...)
			break
		}
	}
	if len(t.children[parentID]) == 0 { // cleanup parent entry if no children
		delete(t.children, parentID)
	}
	delete(t.nodes, id) // remove child from nodes
	return nil
}

// Moves a node.
// Compared to remove and adding a node, this method completes all checks before modifying the tree.
// Ensures that the operation is fully completed.
// Errors if the new parent does not exist.
// Errors if the node is the root or deleted node.
func (t *Tree[MD]) Move(id uuid.UUID, node *TreeNode[MD]) error {
	if _, exists := t.nodes[node.PrntID]; !exists {
		return MissingNodeIDError{id: node.PrntID}
	}
	if id == RootUUID || id == TombstoneUUID {
		return InvalidNodeDeletionError{id: id}
	}

	t.Remove(id)    // remove node from old parent
	t.Add(id, node) // add node to new parent

	return nil
}

// Removes a node and all of its children from the tree.
// Extension to the CRDT algorithm. Allows removal of deleted nodes.
func (t *Tree[MD]) DeleteSubTree(id uuid.UUID) error {
	if _, exists := t.nodes[id]; !exists {
		return MissingNodeIDError{id: id}
	}
	if id == RootUUID || id == TombstoneUUID {
		return InvalidNodeDeletionError{id: id}
	}

	// remove children
	for _, childID := range t.children[id] {
		err := t.DeleteSubTree(childID)
		if err != nil {
			return err
		}
	}

	// remove node
	t.Remove(id)

	delete(t.children, id) // remove child from children
	return nil
}

// Checks if node is ancestor of other node.
// E.g. If ancID can be reached by following the parent pointers from childID.
func (t *Tree[MD]) IsAncestor(childID uuid.UUID, ancID uuid.UUID) (bool, error) {
	if _, exists := t.nodes[childID]; !exists {
		return false, MissingNodeIDError{id: childID}
	}
	if _, exists := t.nodes[ancID]; !exists {
		return false, MissingNodeIDError{id: ancID}
	}
	for childID != RootUUID {
		childID = t.nodes[childID].PrntID
		if bytes.Equal(childID[:], ancID[:]) {
			return true, nil
		}
	}
	return false, nil
}

func (t *Tree[MD]) Contains(id uuid.UUID) bool {
	_, exists := t.nodes[id]
	return exists
}

func (t *Tree[MD]) String() string {
	return fmt.Sprintf("Nodes: %v \nChildren: %v", t.nodes, t.children)
}
