package lumina

import (
	"testing"

	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

func TestLuminaOpInherits(t *testing.T) {
	var op Operation[*c.VectorTimestamp]
	var timestamp *c.VectorTimestamp

	op = &OpRemove[*c.VectorTimestamp]{Timestmp: c.NewVectorTimestamp(), ChldID: uuid.New()}
	timestamp = op.Timestamp()
	timestampRem := op.Timestamp()

	op = &OpAdd[string, *c.VectorTimestamp]{Timestmp: c.NewVectorTimestamp(), ChldID: uuid.New()}
	timestamp = op.Timestamp()
	timestampAdd := op.Timestamp()

	op = &OpMove[string, *c.VectorTimestamp]{Timestmp: c.NewVectorTimestamp(), ChldID: uuid.New()}
	timestamp = op.Timestamp()
	timestampMove := op.Timestamp()

	if timestampRem.Compare(timestamp) != 0 {
		t.Errorf("Timestamps are not equal")
	}

	if timestampAdd.Compare(timestampRem) != 0 && timestampAdd.Compare(timestampMove) != 0 {
		t.Errorf("Timestamps are not equal")
	}
}

func TestLuminaMoveOps(t *testing.T) {
	opMov := &OpMove[string, *c.VectorTimestamp]{Timestmp: c.NewVectorTimestamp(), ChldID: uuid.New(), NewP: treecrdt.NewTreeNode(uuid.New(), "meta"), Priotity: *c.NewLamport()}
	opMov.Timestamp().Inc(uuid.New())

	if opMov.NewP.Meta != "meta" {
		t.Errorf("Wrong meta")
	}

	opMov2 := &OpMove[string, *c.VectorTimestamp]{Timestmp: c.NewVectorTimestamp(), ChldID: uuid.New(), NewP: treecrdt.NewTreeNode(uuid.New(), "meta2"), Priotity: *c.NewLamport()}
	opMov2.Timestmp.Inc(uuid.New())
	opMov2.Priotity.Inc()

	if opMov.CompareOp(opMov2) != 2 {
		t.Errorf("Wrong comparison, expected 2, got %d", opMov.CompareOp(opMov2))
	}

}
