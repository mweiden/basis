package request_count

import (
	"sync"
)

type RequestCount struct {
	timestamp int64
	count     uint64
}

type RequestMetrics struct {
	getTimestamp    func() int64
	interval        int64
	timestampCounts []RequestCount
	counterMutex    sync.RWMutex
}

func NewRequestMetrics(interval int64, getTimestamp func() int64) RequestMetrics {
	if interval < 0 {
		panic(1)
	}
	r := RequestMetrics{}
	r.counterMutex = sync.RWMutex{}
	r.getTimestamp = getTimestamp
	r.interval = interval
	return r
}

func (r *RequestMetrics) sanityCheckTimestamp(timestamp int64) {
	if timestamp < 0 {
		panic(1)
	}
	lastInd := len(r.timestampCounts) - 1
	// time should be monotonically increasing
	if lastInd >= 0 && r.timestampCounts[lastInd].timestamp > timestamp {
		panic(1)
	}
}

func (r *RequestMetrics) Inc(amount uint) {
	if amount != 0 {
		r.counterMutex.Lock()
		r.garbageCollect()
		timestamp := r.getTimestamp()
		r.sanityCheckTimestamp(timestamp)
		lastInd := len(r.timestampCounts) - 1
		overlap := lastInd >= 0 &&
			r.timestampCounts[lastInd].timestamp == timestamp
		if !overlap {
			r.timestampCounts = append(
				r.timestampCounts,
				RequestCount{timestamp, 1},
			)
		} else {
			r.timestampCounts[lastInd].count++
		}
		r.counterMutex.Unlock()
	}
}

func (r *RequestMetrics) Count() uint64 {
	r.counterMutex.RLock()
	end := r.getTimestamp()
	var start int64
	if r.interval > end {
		start = 0
	} else {
		start = end - r.interval
	}
	// pass by the old timestamps out of range, possible cases are:
	// 1. if the next lowest timestamp is not in storage, sum from the first timestamp
	// 2. if the timestamp is found, start from the next one
	i := binarySearch(start, len(r.timestampCounts), r.timestampAt) + 1
	// sum timestamps that are in range
	var total uint64
	nTimestamps := len(r.timestampCounts)
	for i < nTimestamps && r.timestampCounts[i].timestamp <= end {
		total += r.timestampCounts[i].count
		i++
	}
	r.counterMutex.RUnlock()
	return total
}

func (r *RequestMetrics) garbageCollect() {
	timestamp := r.getTimestamp()
	r.sanityCheckTimestamp(timestamp)
	start := timestamp - r.interval
	// remove the timestamp counts before start, there are two cases:
	// 1. if the next lowest timestamp is not in storage, slice from 0
	// 2. if the timestamp is found, slice from the next one
	i := binarySearch(start, len(r.timestampCounts), r.timestampAt) + 1
	r.timestampCounts = r.timestampCounts[i:len(r.timestampCounts)]
}

func (r *RequestMetrics) timestampAt(i int) int64 {
	return r.timestampCounts[i].timestamp
}

// altered binary search
// returns the index of the timestamp or next smallest one if the timestamp DNE
func binarySearch(timestamp int64, length int, getVal func(int) int64) int {
	if length == 0 {
		return -1
	}
	left := 0
	right := length - 1
	for left <= right {
		split := (left + right) / 2
		diff := getVal(split) - timestamp
		if diff == 0 {
			return split
		} else if diff < 0 {
			left = split + 1
		} else {
			right = split - 1
		}
	}
	return left - 1
}
