package datastructures

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"math"
	"sync"
)

type pointValuePair struct {
	point float64
	value interface{}
}

type ConsistentHash struct {
	circlePoints []pointValuePair
	mutex        sync.RWMutex
	replicas     int
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		mutex:    sync.RWMutex{},
		replicas: replicas,
	}
}

func (h *ConsistentHash) Remove(key interface{}) {
	h.mutex.Lock()
	for replica := 0; replica < h.replicas; replica++ {
		point := mapToUnitRange(key, replica)
		h.removePoint(point)
	}
	h.mutex.Unlock()
}

func (h *ConsistentHash) removePoint(point float64) {
	i, found := h.binarySearch(point)
	if found {
		copy(h.circlePoints[i:], h.circlePoints[i+1:])
		h.circlePoints = h.circlePoints[:len(h.circlePoints)-1]
	}
}

func (h *ConsistentHash) Add(key interface{}, value interface{}) {
	h.mutex.Lock()
	for replica := 0; replica < h.replicas; replica++ {
		point := mapToUnitRange(key, replica)
		h.addPoint(point, value)
	}
	h.mutex.Unlock()
}

func (h *ConsistentHash) addPoint(point float64, value interface{}) {
	i, found := h.binarySearch(point)
	pair := pointValuePair{
		point: point,
		value: value,
	}
	if found {
		h.circlePoints[i] = pair
	} else {
		h.circlePoints = append(
			h.circlePoints[:i],
			append([]pointValuePair{pair}, h.circlePoints[i:]...)...,
		)
	}
}

func (h *ConsistentHash) Hash(key interface{}) interface{} {
	h.mutex.RLock()
	keyRangeVal := mapToUnitRange(key, -1)
	i, _ := h.binarySearch(keyRangeVal)
	// wrap to the beginning of the circle if keyRangeValue is greater than the
	// greatest stored value on the circle
	if i > len(h.circlePoints)-1 {
		i = 0
	}
	result := h.circlePoints[i].value
	h.mutex.RUnlock()
	return result
}

// altered binary search
// returns the index of the key or next greatest one if the key DNE
func (h *ConsistentHash) binarySearch(point float64) (int, bool) {
	length := len(h.circlePoints)
	left := 0
	right := length - 1
	for left <= right {
		split := (left + right) / 2
		diff := h.circlePoints[split].point - point
		if diff == 0 {
			return split, true
		} else if diff < 0 {
			left = split + 1
		} else {
			right = split - 1
		}
	}
	return left, false
}

func mapToUnitRange(o interface{}, replica int) float64 {
	h := fnv.New32()
	bytes, err := getBytes(o)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(1)
	}
	h.Write(bytes)
	bytes, err = getBytes(replica)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(1)
	}
	h.Write(bytes)
	return float64(h.Sum32()%math.MaxUint32) / float64(math.MaxUint32)
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
