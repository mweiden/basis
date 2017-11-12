package datastructures

import (
	"reflect"
	"sync"
	"testing"
)

func TestMapToUnitRange(t *testing.T) {
	t.Parallel()
	if mapToUnitRange(1) != mapToUnitRange(1) {
		t.Error("Values should be the same!")
	}
	if mapToUnitRange(1) == mapToUnitRange(2) {
		t.Error("Values should be the different!")
	}
	if mapToUnitRange("a") != mapToUnitRange("a") {
		t.Error("Values should be the different!")
	}
	if mapToUnitRange("a") == mapToUnitRange("b") {
		t.Error("Values should be the different!")
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()
	h := NewConsistentHash()
	expected := &ConsistentHash{
		[]circleKeyValuePair{
			circleKeyValuePair{
				unitVal: mapToUnitRange(10),
				value:   true,
			},
		},
		sync.RWMutex{},
	}
	h.Add(10, true)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("%v != %v", expected, h)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	h := NewConsistentHash()
	expected := NewConsistentHash()
	h.Add(10, true)
	h.Add(20, true)
	h.Add(30, true)
	h.Remove(20)
	expected.Add(10, true)
	expected.Add(30, true)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("%v != %v", expected, h)
	}
}

func TestHash(t *testing.T) {
	t.Parallel()
	h := NewConsistentHash()
	h.Add(10, true)
	h.Add(20, false)
	h.Add(30, false)
	expected := h.Hash(10)
	h.Remove(20)
	result := h.Hash(10)
	if expected != result {
		t.Errorf("%v != %v", expected, result)
	}
}

func TestBinarySearch(t *testing.T) {
	t.Parallel()
	h := ConsistentHash{
		[]circleKeyValuePair{
			circleKeyValuePair{unitVal: 1},
			circleKeyValuePair{unitVal: 2},
			circleKeyValuePair{unitVal: 3},
			circleKeyValuePair{unitVal: 5},
			circleKeyValuePair{unitVal: 6},
			circleKeyValuePair{unitVal: 7},
		},
		sync.RWMutex{},
	}
	expectedResult := 3
	expectedFound := false
	result, found := h.binarySearch(4)
	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
	expectedResult = 3
	expectedFound = true
	result, found = h.binarySearch(5)
	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}

	h = ConsistentHash{
		[]circleKeyValuePair{
			circleKeyValuePair{unitVal: 1},
		},
		sync.RWMutex{},
	}
	expectedResult = 1
	expectedFound = false
	result, found = h.binarySearch(2)
	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}

	h = ConsistentHash{
		[]circleKeyValuePair{},
		sync.RWMutex{},
	}
	expectedResult = 0
	expectedFound = false
	result, found = h.binarySearch(1)
	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
}
