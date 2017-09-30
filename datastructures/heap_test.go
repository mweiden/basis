package datastructures

import (
	"testing"
)

func TestHeap(t *testing.T) {
	t.Parallel()
	h := Heap{
		Compare: func(a interface{}, b interface{}) int {
			return a.(int) - b.(int)
		},
	}
	h.Insert(3)
	h.Insert(2)
	h.Insert(1)
	h.Insert(4)
	type Result struct {
		val interface{}
		err error
	}
	expected := []Result{
		Result{1, nil},
		Result{2, nil},
		Result{3, nil},
		Result{4, nil},
		Result{nil, EOH},
	}

	for _, e := range expected {
		val, err := h.Pop()
		if e.err != err || e.val != val {
			t.Errorf("Expected %v %d, got %v %d", e.err, e.val, err, val)
		}
	}
}
