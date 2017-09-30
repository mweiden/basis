package request_count

import (
	"github.com/basis/datastructures"
	"sync"
)

type RequestMetrics struct {
	getTimestamp    func() uint64
	interval        uint64
	timestampCounts map[uint64]uint64
	timestamps      *datastructures.Heap
	counterMutex    sync.RWMutex
}

func (r *RequestMetrics) Init() {
	r.counterMutex = sync.RWMutex{}
	r.timestampCounts = make(map[uint64]uint64)
	r.timestamps = &datastructures.Heap{
		Compare: func(a interface{}, b interface{}) int {
			return int(a.(uint64) - b.(uint64))
		},
	}
}

func (r *RequestMetrics) Inc(amount uint) {
	if amount != 0 {
		r.counterMutex.Lock()
		r.garbageCollect()
		timestamp := r.getTimestamp()
		_, ok := r.timestampCounts[timestamp]
		if !ok {
			r.timestampCounts[timestamp] = 1
			r.timestamps.Insert(timestamp)
		} else {
			r.timestampCounts[timestamp] += 1
		}
		r.counterMutex.Unlock()
	}
}

func (r *RequestMetrics) Count() uint64 {
	r.counterMutex.RLock()
	end := r.getTimestamp()
	var start uint64
	if (r.interval > end) {
		start = 0
	} else {
		start = end - r.interval
	}
	var total uint64 = 0
	for timestamp, count := range r.timestampCounts {
		if timestamp > start && timestamp <= end {
			total += count
		}
	}
	r.counterMutex.RUnlock()
	return total
}

func (r *RequestMetrics) garbageCollect() {
	end := r.getTimestamp()
	start := end - r.interval
	min, err := r.timestamps.Peek()
	for err == nil && min.(uint64) < start {
		delete(r.timestampCounts, min.(uint64))
		r.timestamps.Pop()
		min, err = r.timestamps.Peek()
	}
}
