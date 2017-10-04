package datastructures

import (
	"testing"
)

func TestLinkedList_Init(t *testing.T) {
	t.Parallel()
	ll := NewLinkedList()
	if ll.start != nil || ll.end != nil {
		t.Error("Not properly initialized.")
	}
}
