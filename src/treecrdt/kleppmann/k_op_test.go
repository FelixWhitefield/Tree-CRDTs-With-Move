package kleppmann

// This will test both ops and logops

import (
	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
	"testing"
)

func TestKleppOpNewAndCompareOp(t *testing.T) {
	u1 := uuid.New()
	l := c.NewLamport()

	op1 := NewOpMove(l.Clone(), RootUUID, u1, "meta")
	op2 := NewOpMove(l.Clone(), RootUUID, u1, "meta")

	if op1.Compare(op2) != 0 {
		t.Errorf("Compare return %d, expected 0", op1.Compare(op2))
	}

	op3 := NewOpMove(l.Tick(), RootUUID, u1, "meta2")

	if op1.Compare(op3) != -1 {
		t.Errorf("Compare return %d, expected -1", op1.Compare(op3))
	}

	// LOG OPS
	lopop := NewLogOpMove(op1, NewTreeNode(uuid.New(), "meta"))
	if lopop.CompareOp(op1) != 0 {
		t.Errorf("Compare return %d, expected 0", lopop.CompareOp(op1))
	}

	lopop2 := NewLogOpMove(op3, NewTreeNode(uuid.New(), "meta2"))
	if lopop2.CompareOp(op1) != 1 {
		t.Errorf("Compare return %d, expected 1", lopop2.CompareOp(op3))
	}
}
