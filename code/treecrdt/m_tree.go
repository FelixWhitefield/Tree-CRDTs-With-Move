// Implements the tree data structure used by the CRDT.
// The `nodes` map represents the triples indexed by child id.
// The `children` map provides a quick lookup of a node's children.
// RootUUID and TombstoneUUID are special nodes that are always present.
// They are set manually as to ensure they are the same across all replicas.
package treecrdt

import (
	"fmt"
	"errors"
	"github.com/google/uuid"
)

var (
	RootUUID      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	TombstoneUUID = uuid.Nil
)

type Tree[MD Metadata] struct {
	nodes    map[uuid.UUID]*TreeNode[MD] // node id -> tree node
	children map[uuid.UUID][]uuid.UUID // node id -> []child id
}

func NewTree[MD Metadata]() *Tree[MD] {
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
func (t *Tree[MD]) GetNode(id uuid.UUID) (*TreeNode[MD], bool) {
	node, exists := t.nodes[id]
	return node, exists
}

func (t *Tree[MD]) GetChildren(id uuid.UUID) ([]uuid.UUID, bool) {
	children, exists := t.children[id]
	return children, exists
}

// Adds a node to the tree.
// Errors if the node already exists.
// Errors if the parent node does not exist.
func (t *Tree[MD]) Add(id uuid.UUID, node *TreeNode[MD]) error {
	if _, exists := t.nodes[id]; exists {
		return errors.New("node already exists")
	}
	if _, exists := t.nodes[node.parentID]; !exists {
		return errors.New("parent node does not exist")
	}

	t.nodes[id] = node
	if _, exists := t.children[node.parentID]; !exists {
		t.children[node.parentID] = []uuid.UUID{id}
	} else {
		t.children[node.parentID] = append(t.children[node.parentID], id)
	}
	return nil
}


// Removes a node from the tree. Doesn't remove the corresponding children entry.
// Errors if the node does not exist.
// Errors if the node is the root or deleted node.
func (t *Tree[MD]) Remove(id uuid.UUID) error {
	if _, exists := t.nodes[id]; !exists {
		return errors.New("node does not exist")
	}
	if id == RootUUID || id == TombstoneUUID {
		return errors.New("cannot remove root or deleted node")
	}

	parentID := t.nodes[id].parentID
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
// Errors if either the node or the new parent does not exist.
// Errors if the node is the root or deleted node.
func (t *Tree[MD]) Move(id uuid.UUID, node *TreeNode[MD]) error {
	if _, exists := t.nodes[id]; !exists {
		return errors.New("node does not exist")
	}
	if _, exists := t.nodes[node.parentID]; !exists {
		return errors.New("new parent node does not exist")
	}
	if id == RootUUID || id == TombstoneUUID {
		return errors.New("cannot move root or deleted node")
	}

	t.Remove(id) // remove node from old parent
	t.Add(id, node) // add node to new parent

	return nil
}

// Removes a node and all of its children from the tree.
// Extension to the CRDT algorithm. Allows removal of deleted nodes.
func (t *Tree[MD]) DeleteSubTree(id uuid.UUID) error {
	if _, exists := t.nodes[id]; !exists {
		return errors.New("node does not exist")
	}
	if id == RootUUID || id == TombstoneUUID {
		return errors.New("cannot remove root or deleted node")
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
func (t *Tree[MD]) IsAncestor(childID uuid.UUID, ancID uuid.UUID) bool {
	for childID != RootUUID {
		if childID = t.nodes[childID].parentID; childID == ancID {
			return true
		}
	}
	return false
}

func (t *Tree[MD]) String() string {
	return fmt.Sprintf("nodes: %v \nchildren: %v", t.nodes, t.children)
}
