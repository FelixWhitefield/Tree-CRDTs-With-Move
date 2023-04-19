package treecrdt_test

import (
	"testing"

	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	u "github.com/google/uuid"
)

func TestApplyOp(t *testing.T) {
	tp := NewTreeReplica[string]()

	uid := u.New();
	op := tp.Prepare(uid, RootUUID, "hello")

	tp.Effect(*op)
	tp.Effect(*op)

	child, _:= tp.GetChildren(RootUUID)
	if child[0] != uid {
		t.Errorf("Invalid ID child of root %v, should be %v", child[0], uid)
	}

	id2 := u.New() 
	op2 := tp.Prepare(id2, uid, "2nd")

	tp.Effect(*op2)

	child, _ = tp.GetChildren(uid)
	if child[0] != id2 {
		t.Errorf("Invalid ID child of root %v, should be %v", child[0], uid)
	}

	op4 := tp.Prepare(uid, id2, "hello") 
	tp.Effect(*op4)

	child, _ = tp.GetChildren(id2)
	if len(child) > 0 {
		t.Errorf("There is a cycle")
	}

	// op3 := tp.Prepare(id2, TombstoneUUID, "2nd")
	// tp.Effect(*op3)

	// child, _ = tp.GetChildren(uid)
	// if len(child) != 0 {
	// 	t.Errorf("Error deleting node")
	// }
}