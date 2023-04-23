package treecrdt_test

import (
	"testing"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/k"
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	u "github.com/google/uuid"
)

func TestNewAndAdd(t *testing.T) {
	uuid1 := u.New()
	tr := NewTreeReplica[string](nil, uuid1)

	if (tr == nil) {
		t.Errorf("NewTreeReplicaWithID() returned nil")
	}

	if (tr != nil && tr.ActorID() != uuid1) {
		t.Errorf("NewTreeReplicaWithID() returned wrong ID")
	}

	uuid2 := u.New()
	op := tr.Prepare(uuid2, RootUUID, "test")
	tr.Effect(op)

	if c, _ := tr.GetChildren(RootUUID); c[0] != uuid2 {
		t.Errorf("Effect() did not add node correctly")
	}

	if n := tr.GetNode(uuid2); n.ParentID() != RootUUID {
		t.Errorf("Effect() did not add node correctly")
	}

	if n := tr.GetNode(uuid2); n.Metadata() != "test" {
		t.Errorf("Effect() did not add node correctly")
	}

	uuid3 := u.New()
	op = tr.Prepare(uuid3, uuid2, "test2")
	tr.Effect(op)

	if c, _ := tr.GetChildren(RootUUID); contains(c, uuid2) && contains(c, uuid3) { 
		t.Errorf("Effect() did not add node correctly")
	}
}

func contains(s []u.UUID, e u.UUID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestTime(t *testing.T) {
	tr := NewTreeReplica[string](nil)

	uuid1 := u.New()
	op := tr.Prepare(uuid1, RootUUID, "test")
	tr.Effect(op)

	expectedTime := c.NewLamport(tr.ActorID())
	expectedTime.Inc()
	
	if (tr.CurrentTime().ActorID() != uuid1 &&  tr.CurrentTime().Compare(expectedTime) != 0) {
		t.Errorf("CurrentTime() returned wrong time")
	}
}