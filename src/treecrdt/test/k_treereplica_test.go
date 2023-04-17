package treecrdt_test

import (
	"fmt"
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
	fmt.Println(child)
	if child[0] != uid {
		t.Errorf("Invalid ID child of root %v, should be %v", child[0], uid)
	}
}