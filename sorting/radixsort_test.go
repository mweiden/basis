package sorting

import (
	"reflect"
	"testing"
)

func TestRadixSort(t *testing.T) {
	t.Parallel()
	toSort := []int{20, 35, 15, 1, 43, 67, 126, 512}
	expected := []int{1, 15, 20, 35, 43, 67, 126, 512}
	RadixSort(toSort)
	if !reflect.DeepEqual(toSort, expected) {
		t.Errorf("%v != %v\n", toSort, expected)
	}
}
