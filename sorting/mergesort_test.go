package sorting

import (
	"reflect"
	"testing"
)

func TestMergeSort(t *testing.T) {
	t.Parallel()
	trail := []int{2, 1, 0, 3}
	MergeSort(trail)

	expected := []int{0, 1, 2, 3}
	if !reflect.DeepEqual(trail, expected) {
		t.Errorf("%v != %v!", trail, expected)
	}
}
