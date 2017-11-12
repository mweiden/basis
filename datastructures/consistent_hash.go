package datastructures

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"math"
	"sync"
)

type circleKeyValuePair struct {
	unitVal float64
	value   interface{}
}

type ConsistentHash struct {
	circlePoints []circleKeyValuePair
	mutex        sync.RWMutex
}

func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{mutex: sync.RWMutex{}}
}

func (h *ConsistentHash) Remove(key interface{}) {
	h.mutex.Lock()
	keyRangeVal := mapToUnitRange(key)
	i, found := h.binarySearch(keyRangeVal)
	if found {
		tmp := h.circlePoints
		h.circlePoints = h.circlePoints[0:i]
		if i+1 < len(h.circlePoints) {
			h.circlePoints = append(h.circlePoints, tmp[i+1:len(tmp)]...)
		}
	}
	h.mutex.Unlock()
}

func (h *ConsistentHash) Add(key interface{}, value interface{}) {
	h.mutex.Lock()
	keyRangeVal := mapToUnitRange(key)
	i, found := h.binarySearch(keyRangeVal)
	kvPair := circleKeyValuePair{
		unitVal: keyRangeVal,
		value:   value,
	}
	if found {
		h.circlePoints[i] = kvPair
	} else {
		var newCirclePoints []circleKeyValuePair
		for _, p := range h.circlePoints[0:i] {
			newCirclePoints = append(newCirclePoints, p)
		}
		newCirclePoints = append(newCirclePoints, kvPair)
		if len(h.circlePoints) > 0 {
			for _, p := range h.circlePoints[i:len(h.circlePoints)] {
				newCirclePoints = append(newCirclePoints, p)
			}
		}
		h.circlePoints = newCirclePoints
	}
	h.mutex.Unlock()
}

func (h *ConsistentHash) Hash(key interface{}) interface{} {
	h.mutex.RLock()
	keyRangeVal := mapToUnitRange(key)
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
func (h *ConsistentHash) binarySearch(searchVal float64) (int, bool) {
	length := len(h.circlePoints)
	left := 0
	right := length - 1
	for left <= right {
		split := (left + right) / 2
		diff := h.circlePoints[split].unitVal - searchVal
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

func mapToUnitRange(o interface{}) float64 {
	h := fnv.New32()
	bytes, err := getBytes(o)
	if err != nil {
		panic(1)
	}
	h.Write(bytes)
	return float64(h.Sum32()) / float64(math.MaxUint32)
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
