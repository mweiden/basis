package request_count

import (
	"sync"
)

type RequestMetrics struct {
	getTimestamp func() uint64
	interval uint64
	maxSize uint64
	timestampCounts map[uint64]uint64
	counterMutex *sync.Mutex
}

func (r *RequestMetrics) Init() {
	r.counterMutex = &sync.Mutex{}
	r.maxSize = r.interval * 2
	r.timestampCounts = make(map[uint64]uint64)
}

func (r *RequestMetrics) Inc(amount int) {
	// clean up if the size is getting out of hand
	if uint64(len(r.timestampCounts)) > r.maxSize {
		r.garbageCollect()
	}
	// increment
	r.counterMutex.Lock()
	timestamp := r.getTimestamp()
	_, ok := r.timestampCounts[timestamp]
	if !ok {
		r.timestampCounts[timestamp] = 1
	} else {
		r.timestampCounts[timestamp] += 1
	}
	r.counterMutex.Unlock()
}

func (r *RequestMetrics) Count() uint64 {
	r.counterMutex.Lock()
	end := r.getTimestamp()
	start := end - r.interval
	var total uint64 = 0
	for timestamp, count := range r.timestampCounts {
		if timestamp > start && timestamp <= end {
			total += count
		}
	}
	r.counterMutex.Unlock()
	return total
}

func (r *RequestMetrics) garbageCollect() uint64 {
	r.counterMutex.Lock()
	end := r.getTimestamp()
	start := end - r.interval
	var discardCount uint64 = 0

	newMap := make(map[uint64]uint64)
	for timestamp, count := range r.timestampCounts {
		if timestamp >= start && timestamp < end {
			newMap[timestamp] = count
		} else {
			discardCount += count
		}
	}
	r.timestampCounts = newMap
	r.counterMutex.Unlock()
	return discardCount
}
