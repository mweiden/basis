package request_count

import (
	"sync"
	"testing"
)

func TestRequestMetrics(t *testing.T) {
	t.Parallel()

	var timestamp int64
	getTimestamp := func() int64 {
		return timestamp
	}
	metrics := NewRequestMetrics(5000, getTimestamp)
	val := metrics.Count()
	var expected uint64
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}

	timestamp = 10000
	metrics.Inc(1)
	val = metrics.Count()
	expected = 1
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}

	timestamp = 10001
	metrics.Inc(1)
	val = metrics.Count()
	expected = 2
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}

	timestamp = 14999
	val = metrics.Count()
	expected = 2
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}

	timestamp = 15000
	val = metrics.Count()
	expected = 1
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}

	timestamp = 15001
	val = metrics.Count()
	expected = 0
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}
}

func TestRequestMetrics_Inc(t *testing.T) {
	t.Parallel()
	var timestamp int64 = 10000
	getTimestamp := func() int64 { return timestamp }
	metrics := NewRequestMetrics(5000, getTimestamp)

	// should lock write path properly
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			metrics.Inc(1)
			wg.Done()
		}()
	}
	wg.Wait()

	var expected uint64 = 50
	val := metrics.Count()

	if val != expected {
		t.Errorf("%v != %v\n", val, expected)
	}
}

func TestRequestMetrics_garbageCollect(t *testing.T) {
	t.Parallel()
	var timestamp int64
	getTimestamp := func() int64 {
		return timestamp
	}
	metrics := NewRequestMetrics(5000, getTimestamp)
	for i := 0; i < 10; i++ {
		timestamp = 10000 + int64(i)
		metrics.Inc(1)
	}
	for i := 0; i < 10; i++ {
		timestamp = 15000 + int64(i)
		metrics.garbageCollect()
		expected := 9 - i
		mapSize := len(metrics.timestampCounts)
		if mapSize != expected {
			t.Errorf("iter=%d, %d != %d", i, mapSize, expected)
		}
	}
}

func TestBinarySearch(t *testing.T) {
	t.Parallel()
	ary := []int64{1, 2, 3, 5, 6, 7}
	getVal := func(i int) int64 {
		return ary[i]
	}
	expected := 2
	result := binarySearch(4, len(ary), getVal)
	if expected != result {
		t.Errorf("%d != %d", expected, result)
	}
	expected = 3
	result = binarySearch(5, len(ary), getVal)
	if expected != result {
		t.Errorf("%d != %d", expected, result)
	}
	ary = []int64{1}
	getVal = func(i int) int64 {
		return ary[i]
	}
	expected = 0
	result = binarySearch(2, len(ary), getVal)
	if expected != result {
		t.Errorf("%d != %d", expected, result)
	}
	ary = []int64{}
	getVal = func(i int) int64 {
		return ary[i]
	}
	expected = -1
	result = binarySearch(1, len(ary), getVal)
	if expected != result {
		t.Errorf("%d != %d", expected, result)
	}
}
