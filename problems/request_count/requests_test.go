package request_count

import (
	"testing"
	"sync"
)

func TestRequestMetrics(t *testing.T) {
	t.Parallel()

	var timestamp uint64
	getTimestamp := func() uint64 {
		return timestamp
	}
	metrics := RequestMetrics{
		getTimestamp: getTimestamp,
		interval:     5000,
	}

	metrics.Init()
	val := metrics.Count()
	var expected uint64 = 0
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
	var timestamp uint64 = 10000
	getTimestamp := func() uint64 { return timestamp }
	metrics := RequestMetrics{
		getTimestamp: getTimestamp,
		interval:     5000,
	}
	metrics.Init()

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
	var timestamp uint64
	getTimestamp := func() uint64 {
		return timestamp
	}
	metrics := RequestMetrics{
		getTimestamp: getTimestamp,
		interval:     5000,
	}
	metrics.Init()
	for i := 0; i < 10; i++ {
		timestamp = 10000 + uint64(i)
		metrics.Inc(1)
	}
	for i := 0; i < 10; i++ {
		timestamp = 15000 + uint64(i)
		metrics.garbageCollect()
		expected := 10 - i
		arySize := len(metrics.timestamps)
		mapSize := len(metrics.timestampCounts)
		if arySize != mapSize {
			t.Errorf("iter=%d, %d != %d", i, arySize, mapSize)
		}
		if arySize != expected {
			t.Errorf("iter=%d, %d != %d", i, expected, arySize)
		}
	}
}
