package main

import (
	"testing"
)

func TestLinkedList_Init(t *testing.T) {
	t.Parallel()
	var ll LinkedList
	ll.Init()
	if ll.start != nil || ll.end != nil {
		t.Error("Not properly initialized.")
	}
}
