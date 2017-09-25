package main

import (
	"testing"
)

func TestStack(t *testing.T) {
	t.Parallel()
	var s Stack
	s.Push(1)
	s.Push(2)

	expected := []int{2, 1}
	for _, i := range expected {
		err, result := s.Pop()
		if result != i || err != nil {
			t.Errorf("Expected (%v, %d), got (%v, %d)!", nil, i, err, result)
		}
	}

	err, result := s.Pop()
	if err != EOS || result != -1 {
		t.Errorf("Expected (%v, %d), got (%v, %d)!", nil, -1, err, result)
	}
}
