package hash

import (
	"reflect"
	"sync"
	"testing"
)

func TestMapToUnitRange(t *testing.T) {
	t.Parallel()
	if mapToUnitRange(1, 0) != mapToUnitRange(1, 0) {
		t.Error("Values should be the same!")
	}
	if mapToUnitRange(1, 0) == mapToUnitRange(2, 0) {
		t.Error("Values should be the different!")
	}
	if mapToUnitRange(1, 0) == mapToUnitRange(1, 1) {
		t.Error("Values should be the different!")
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()
	h := NewConsistentHash(1)
	expected := &ConsistentHash{
		circlePoints: []pointValuePair{
			pointValuePair{
				point: mapToUnitRange(10, 0),
				value: true,
			},
		},
		mutex:    sync.RWMutex{},
		replicas: 1,
	}
	h.Add(10, true)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("%v != %v", expected, h)
	}
	// adding repeatedly should not change the datastructure
	h.Add(10, true)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("%v != %v", expected, h)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	// with one replica
	h := NewConsistentHash(1)
	expected := NewConsistentHash(1)
	h.Add(10, true)
	h.Add(20, true)
	h.Add(30, true)
	h.Remove(20)
	expected.Add(10, true)
	expected.Add(30, true)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("%v != %v", expected, h)
	}
	// with multiple replicas
	h = NewConsistentHash(2)
	expected = NewConsistentHash(2)
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
	// the hashes should stay consistent after deletion
	h := NewConsistentHash(1)
	h.Add(10, 1.0)
	h.Add(20, 2.0)
	h.Add(30, 3.0)
	expected := h.Hash(10)
	h.Remove(20)
	result := h.Hash(10)
	if expected != result {
		t.Errorf("%v != %v", expected, result)
	}
	h.Remove(10)
	h.Add(10, 1.0)
	result = h.Hash(10)
	if expected != result {
		t.Errorf("%v != %v", expected, result)
	}
	// the hash should wrap in the point circle
	lowPair := h.circlePoints[0]
	highPair := h.circlePoints[len(h.circlePoints)-1]
	if lowPair.point >= highPair.point {
		t.Errorf("low=%v, high=%v!", lowPair, highPair)
	}
	highSearch := mapToUnitRange(22, -1)
	if highSearch <= highPair.point {
		t.Errorf("highSearch=%v, high=%v!", highSearch, highPair)
	}
	highSearchResult := h.Hash(22)
	if highSearchResult != lowPair.value {
		t.Errorf("%v != %v", lowPair.value, highSearchResult)
	}
}
