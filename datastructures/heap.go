package datastructures

import (
	"errors"
)

var (
	EOH = errors.New("End of Heap!")
)

type Heap struct {
	ary     []interface{}
	Compare func(a interface{}, b interface{}) int
}

func (h *Heap) Insert(val interface{}) {
	// append to end and then bubble up
	h.ary = append(h.ary, val)
	i := len(h.ary) - 1
	for i > 0 {
		parentInd := heapParent(i)
		if h.Compare(h.ary[parentInd], h.ary[i]) <= 0 {
			break
		} else {
			swap(h.ary, parentInd, i)
		}
		i = parentInd
	}
}

func (h *Heap) Pop() (interface{}, error) {
	if len(h.ary) == 0 {
		return nil, EOH
	}
	// pull from front, replace with last, then bubble down
	root := h.ary[0]
	h.ary[0] = h.ary[len(h.ary)-1]
	h.ary = h.ary[:len(h.ary)-1]
	i := 0
	for i < len(h.ary)-1 {
		compareInd := i
		l := heapLeft(i)
		r := heapRight(i)
		if l < len(h.ary) && h.Compare(h.ary[l], h.ary[compareInd]) < 0 {
			compareInd = l
		}
		if r < len(h.ary) && h.Compare(h.ary[r], h.ary[compareInd]) < 0 {
			compareInd = r
		}
		if compareInd != i {
			swap(h.ary, compareInd, i)
			i = compareInd
		} else {
			break
		}
	}
	return root, nil
}

func (h *Heap) Peek() (interface{}, error) {
	if len(h.ary) == 0 {
		return nil, EOH
	} else {
		return h.ary[0], nil
	}
}

func (h *Heap) Size() int {
	return len(h.ary)
}

func swap(ary []interface{}, i int, j int) {
	ival := ary[i]
	ary[i] = ary[j]
	ary[j] = ival
}

func heapLeft(i int) int {
	return i*2 + 1
}

func heapRight(i int) int {
	return i*2 + 2
}

func heapParent(i int) int {
	return (i - 1) / 2
}
