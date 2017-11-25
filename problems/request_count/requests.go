package request_count

import (
	"math"
	"sync"

	"github.com/mweiden/basis/search"
)

type RequestCount struct {
	timestamp int64
	count     uint64
}

type RequestMetrics struct {
	getTimestamp    func() int64
	interval        uint32
	timestampCounts []RequestCount
	counterMutex    sync.RWMutex
}

func NewRequestMetrics(interval uint32, getTimestamp func() int64) RequestMetrics {
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
	if amount == 0 {
		return
	}
	timestamp := r.getTimestamp()
	r.sanityCheckTimestamp(timestamp)

	// if storage is growing too large, garbage collect
	if uint32(len(r.timestampCounts)) > r.interval {
		r.garbageCollect()
	}

	// ok, actually increment
	r.counterMutex.Lock()
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

func (r *RequestMetrics) Count() uint64 {
	r.counterMutex.RLock()
	end := r.getTimestamp()
	var start int64
	if int64(r.interval) > end {
		start = 0
	} else {
		start = end - int64(r.interval)
	}
	nTimestamps := len(r.timestampCounts)
	i := 0
	found := false
	// if start > timestampAt(splitInd), there are more than log2(n)
	// operations to determine the start of timestamps within ttl, so
	// it's worth it to do binary search
	splitInd := int(math.Log2(float64(nTimestamps)))
	if splitInd > 0 && start > r.timestampCounts[splitInd].timestamp {
		// pass by the old timestamps out of range, possible cases are:
		// 1. if start is less than the first timestamp in storage, sum from the first timestamp
		// 2. if the timestamp is found, sum from the next one
		// 3. if a lesser timestamp is found, sum from the next one
		i, found = By(
			func(i int, counts []RequestCount) int64 {
				return counts[i].timestamp
			},
		).Search(r.timestampCounts, start)
		if found {
			i++
		}
	} else {
		for i < nTimestamps && r.timestampCounts[i].timestamp <= start {
			i++
		}
	}
	// sum timestamps that are in range
	var total uint64
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
	r.counterMutex.Lock()
	start := timestamp - int64(r.interval)
	// remove the timestamp counts before start, possible cases are:
	// 1. if start is less than the first timestamp in storage, slice from the first timestamp
	// 2. if the timestamp is found, slice from the next one
	// 3. if a lesser timestamp is found, slice from the next one
	i, found := By(
		func(i int, counts []RequestCount) int64 {
			return counts[i].timestamp
		},
	).Search(r.timestampCounts, start)
	if found {
		i++
	}
	r.timestampCounts = r.timestampCounts[i:len(r.timestampCounts)]
	r.counterMutex.Unlock()
}

// searching
type By func(int, []RequestCount) int64

type timestampSearcher struct {
	slice []RequestCount
	by    By
}

func (s *timestampSearcher) Len() int {
	return len(s.slice)
}

func (s *timestampSearcher) Compare(i int, point interface{}) int {
	diff := s.by(i, s.slice) - point.(int64)
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	} else {
		return 0
	}
}

func (by By) Search(slice []RequestCount, point int64) (int, bool) {
	ps := &timestampSearcher{
		slice: slice,
		by:    by,
	}
	return search.BinarySearch(ps, point)
}
