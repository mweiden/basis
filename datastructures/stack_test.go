package datastructures

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
		result, err := s.Pop()
		if result != i || err != nil {
			t.Errorf("Expected (%v, %d), got (%v, %d)!", nil, i, err, result)
		}
	}

	result, err := s.Pop()
	expectedErr := EOS
	var expectedResult interface{} = nil
	if expectedErr != err {
		t.Errorf("%v != %v", expectedErr, err)
	}
	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}
}
