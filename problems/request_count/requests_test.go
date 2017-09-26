package request_count

import (
	"testing"
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

	for i := 0; i < 50; i++ {
		timestamp = 4960 + uint64(i)
		metrics.Inc(1)
	}

	timestamp = 10000
	var expected uint64 = 40
	val := metrics.garbageCollect()
	if val != expected {
		t.Errorf("%d != %d", val, expected)
	}
}
