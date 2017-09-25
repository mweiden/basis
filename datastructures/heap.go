package main

import (
	"errors"
)

var (
	EOH = errors.New("End of Heap!")
)

type Heap struct {
	ary []int
}

func (h *Heap) Insert(val int) {
	// append to end and then bubble up
	h.ary = append(h.ary, val)
	i := len(h.ary) - 1
	for i > 0 {
		parentInd := parent(i)
		if h.ary[parentInd] <= h.ary[i] {
			break
		} else {
			swap(h.ary, parentInd, i)
		}
		i = parentInd
	}
}

func (h *Heap) Pop() (error, int) {
	if len(h.ary) == 0 {
		return EOH, -1
	}
	// pull from front, replace with last, then bubble down
	min := h.ary[0]
	h.ary[0] = h.ary[len(h.ary)-1]
	h.ary = h.ary[:len(h.ary)-1]
	i := 0
	for i < len(h.ary)-1 {
		minInd := i
		l := left(i)
		r := right(i)
		if l < len(h.ary) && h.ary[l] < h.ary[minInd] {
			minInd = l
		}
		if r < len(h.ary) && h.ary[r] < h.ary[minInd] {
			minInd = r
		}
		if minInd != i {
			swap(h.ary, minInd, i)
			i = minInd
		} else {
			break
		}
	}
	return nil, min
}

func swap(ary []int, i int, j int) {
	iVal := ary[i]
	ary[i] = ary[j]
	ary[j] = iVal
}

func left(i int) int {
	return i*2 + 1
}

func right(i int) int {
	return i*2 + 2
}

func parent(i int) int {
	return (i-1) / 2
}
