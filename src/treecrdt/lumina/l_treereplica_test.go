package lumina

import (
	"testing"

	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	//"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	u "github.com/google/uuid"
)

func TestLuminaTreeReplicaPrepareAndEffect(t *testing.T) {
	uuid1 := u.New()
	tr := NewTreeReplica[string](uuid1)

	if tr == nil {
		t.Errorf("NewTreeReplicaWithID() returned nil")
	}

	if tr != nil && tr.ActorID() != uuid1 {
		t.Errorf("NewTreeReplicaWithID() returned wrong ID")
	}

	uuid2 := u.New()
	op := tr.PrepareAdd(uuid2, RootUUID, "test")
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
	op = tr.PrepareAdd(uuid3, uuid2, "test2")
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

func TestLuminaTreeReplicaTime(t *testing.T) {
	tr := NewTreeReplica[string]()

	uuid1 := u.New()
	op := tr.PrepareAdd(uuid1, RootUUID, "test")
	tr.Effect(op)

	expectedTime := c.NewVectorTimestamp()
	expectedTime.Inc(tr.ActorID())

	if tr.CurrentTime().Compare(expectedTime) != 0 {
		t.Errorf("CurrentTime() returned wrong time")
	}
}
