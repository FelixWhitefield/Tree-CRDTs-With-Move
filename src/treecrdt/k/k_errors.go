package k

// Represents the errors which can occur when performing operations on the tree
//
// Usually errors can be ignored as they are part of normal operation

import (
	"fmt"
	"github.com/google/uuid"
)

type MissingNodeIDError struct {
	id uuid.UUID
}

func (e MissingNodeIDError) Error() string {
	return fmt.Sprintf("id does not exist: %s", e.id)
}

type IDAlreadyExistsError struct {
	id uuid.UUID
}

func (e IDAlreadyExistsError) Error() string {
	return fmt.Sprintf("id already exists: %s", e.id)
}

type InvalidNodeDeletionError struct {
	id uuid.UUID
}

func (e InvalidNodeDeletionError) Error() string {
	if e.id == RootUUID {
		return "cannot delete root node"
	} else if e.id == TombstoneUUID {
		return "cannot delete tombstone node"
	}
	return fmt.Sprintf("cannot delete node: %s", e.id)
}
