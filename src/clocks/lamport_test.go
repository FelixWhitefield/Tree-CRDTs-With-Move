package clocks

import (
	"bytes"
	u "github.com/google/uuid"
	"testing"
)

func TestNewAndCompare(t *testing.T) {
	u1 := u.New()
	l1 := NewLamport(u1)
	l2 := NewLamport()

	l1a := l1.ActorID()
	l2a := l2.ActorID()

	if l1a != u1 {
		t.Errorf("NewLamport() returned wrong ID")
	}

	if l1.Compare(l2) != bytes.Compare(l1a[:], l2a[:]) {
		t.Errorf("Clocks not equal, expected %d, got %d", bytes.Compare(l1a[:], l2a[:]), l1.Compare(l2))
	}

	if l1.Compare(l2) == 0 {
		t.Errorf("Clocks equal, expected not equal")
	}

	l1.Inc()

	if l1.Compare(l2) != 1 {
		t.Errorf("Error in Inc() or Compare(), expected 1, got %d", l1.Compare(l2))
	}

	if l2.Compare(l1) != -1 {
		t.Errorf("Error in Inc() or Compare(), expected -1, got %d", l2.Compare(l1))
	}
}

func TestTickAndClone(t *testing.T) {
	l1 := NewLamport()

	l1copy := l1.Clone()

	l1tick := l1.Tick() // Should have no effect

	if l1.Compare(l1copy) != 0 {
		t.Errorf("Error in Tick() or Clone(), expected 0, got %d", l1.Compare(l1copy))
	}

	if l1tick.Compare(l1) != 1 {
		t.Errorf("Error in Tick() or Compare(), expected 1, got %d", l1tick.Compare(l1))
	}
}

func TestMerge(t *testing.T) {
	l1 := NewLamport()
	l2 := NewLamport()

	l1.Inc()
	l2.Inc()

	l1copy := l1.Clone() // make copy
	l1.Merge(l2)         // merge l2 into l1, should have no effect as both counters should be equal

	if l1.Compare(l1copy) != 0 {
		t.Errorf("Error in Merge(), expected 0, got %d", l1.Compare(l2))
	}

	l2.Inc()
	l1.Merge(l2)

	if l1.Compare(l1copy) != 1 {
		t.Errorf("Error in Merge(), expected 1, got %d", l1.Compare(l1copy))
	}
}
