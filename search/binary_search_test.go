package search

import "testing"

type searchable struct {
	ary []int
}

func (s *searchable) Len() int {
	return len(s.ary)
}

func (s *searchable) Compare(i int, point interface{}) int {
	return s.ary[i] - point.(int)
}

func TestBinarySearch(t *testing.T) {
	t.Parallel()
	ary := &searchable{[]int{1, 2, 3, 5, 6, 7}}
	expectedInd := 3
	expectedFound := false
	ind, found := BinarySearch(ary, 4)
	if expectedInd != ind {
		t.Errorf("%d != %d", expectedInd, ind)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
	expectedInd = 3
	expectedFound = true
	ind, found = BinarySearch(ary, 5)
	if expectedInd != ind {
		t.Errorf("%d != %d", expectedInd, ind)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
	ary = &searchable{[]int{1}}
	expectedInd = 1
	expectedFound = false
	ind, found = BinarySearch(ary, 2)
	if expectedInd != ind {
		t.Errorf("%d != %d", expectedInd, ind)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
	ary = &searchable{[]int{}}
	expectedInd = 0
	expectedFound = false
	ind, found = BinarySearch(ary, 1)
	if expectedInd != ind {
		t.Errorf("%d != %d", expectedInd, ind)
	}
	if expectedFound != found {
		t.Errorf("%v != %v", expectedFound, found)
	}
}
