package request_count

import (
	"sync"
)

type RequestMetrics struct {
	getTimestamp    func() uint64
	interval        uint64
	timestampCounts map[uint64]uint64
	timestamps      []uint64
	counterMutex    sync.RWMutex
}

func (r *RequestMetrics) Init() {
	r.counterMutex = sync.RWMutex{}
	r.timestampCounts = make(map[uint64]uint64)
}

func (r *RequestMetrics) Inc(amount uint) {
	if amount != 0 {
		r.counterMutex.Lock()
		r.garbageCollect()
		timestamp := r.getTimestamp()
		_, ok := r.timestampCounts[timestamp]
		if !ok {
			r.timestampCounts[timestamp] = 1
			r.timestamps = append(r.timestamps, timestamp)
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
	if r.interval > end {
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
	start := r.getTimestamp() - r.interval
	for len(r.timestamps) > 0 && r.timestamps[0] < start {
		delete(r.timestampCounts, r.timestamps[0])
		r.timestamps = r.timestamps[1:len(r.timestamps)]
	}
}
