package main

import (
	"testing"
)

func TestHeap(t *testing.T) {
	t.Parallel()
	var h Heap
	h.Insert(3)
	h.Insert(2)
	h.Insert(1)
	h.Insert(4)
	type Result struct {
		err error
		val int
	}
	expected := []Result{
		Result{nil, 1},
		Result{nil, 2},
		Result{nil, 3},
		Result{nil, 4},
		Result{EOH, -1},
	}

	for _, e := range expected {
		err, val := h.Pop()
		if e.err != err || e.val != val {
			t.Errorf("Expected %v %d, got %v %d", e.err, e.val, err, val)
		}
	}
}
