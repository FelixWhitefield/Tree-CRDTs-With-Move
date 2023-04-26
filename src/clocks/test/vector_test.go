package clocks_test

import (
	. "github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	u "github.com/google/uuid"
	"testing"
)

func TestVNewAndCompare(t *testing.T) {
	u1 := u.New()
	v1 := NewVectorClock(u1)
	v2 := NewVectorClock()

	v1a := v1.ActorID()

	if v1a != u1 {
		t.Errorf("NewVectorClock() returned wrong ID")
	}

	if v1.CompareTimestamp(v2.Timestamp()) != 0 {
		t.Errorf("Clocks not equal, expected 0, got %d", v1.CompareTimestamp(v2.Timestamp()))
	}

	v1.Inc()

	if v1.CompareTimestamp(v2.Timestamp()) != 1 {
		t.Errorf("Error in Inc() or Compare(), expected 1, got %d", v1.CompareTimestamp(v2.Timestamp()))
	}

	if v2.CompareTimestamp(v1.Timestamp()) != -1 {
		t.Errorf("Error in Inc() or Compare(), expected -1, got %d", v2.CompareTimestamp(v1.Timestamp()))
	}
}

func TestVTickAndClone(t *testing.T) {
	v1 := NewVectorClock()

	v1copy := v1.Timestamp() // Get a copy of the timestamp
	v1tick := v1.Tick()      // Should have no effect

	if v1.CompareTimestamp(v1copy) != 0 {
		t.Errorf("Error in Tick() or Clone(), expected 0, got %d", v1.CompareTimestamp(v1copy))
	}

	if v1tick.Compare(v1.Timestamp()) != 1 {
		t.Errorf("Error in Tick() or Compare(), expected 1, got %d", v1tick.Compare(v1.Timestamp()))
	}
}

func TestVMerge(t *testing.T) {
	v1 := NewVectorClock()
	v2 := NewVectorClock()

	v1.Inc()
	v2.Inc()

	if v1.CompareTimestamp(v2.Timestamp()) != 2 {
		t.Errorf("Error in Inc() or Compare(), expected 2, got %d", v1.CompareTimestamp(v2.Timestamp()))
	}

	v1.Merge(v2.Timestamp())

	if v1.CompareTimestamp(v2.Timestamp()) != 1 {
		t.Errorf("Error in Merge() or Compare(), expected 0, got %d", v1.CompareTimestamp(v2.Timestamp()))
	}

	if v2.CompareTimestamp(v1.Timestamp()) != -1 {
		t.Errorf("Error in Merge() or Compare(), expected -1, got %d", v2.CompareTimestamp(v1.Timestamp()))
	}

	// Make them equal again
	v2.Merge(v1.Timestamp())

	if v1.CompareTimestamp(v2.Timestamp()) != 0 {
		t.Errorf("Error in Merge() or Compare(), expected 0, got %d", v1.CompareTimestamp(v2.Timestamp()))
	}
}
