package sorting

import (
	"image/color"
	"github.com/basis/datastructures"
)

type InstrumentedHeap struct {
	ary []int
}

func (h *InstrumentedHeap) Insert(val int, max int) [][]color.RGBA {
	var img [][]color.RGBA
	// append to end and then bubble up
	h.ary = append(h.ary, val)
	i := len(h.ary) - 1
	img = append(img, arrayToColors(h.ary, max))
	for i > 0 {
		parentInd := datastructures.HeapParent(i)
		if h.ary[parentInd] <= h.ary[i] {
			break
		} else {
			datastructures.Swap(h.ary, parentInd, i)
		}
		i = parentInd
		img = append(img, arrayToColors(h.ary, max))
	}
	return img
}

func (h *InstrumentedHeap) Pop(max int) (error, int, [][]color.RGBA) {
	var img [][]color.RGBA
	if len(h.ary) == 0 {
		return datastructures.EOH, -1, img
	}
	// pull from front, replace with last, then bubble down
	min := h.ary[0]
	h.ary[0] = h.ary[len(h.ary)-1]
	h.ary = h.ary[:len(h.ary)-1]
	i := 0
	for i < len(h.ary)-1 {
		img = append(img, arrayToColors(h.ary, max))
		minInd := i
		l := datastructures.HeapLeft(i)
		r := datastructures.HeapRight(i)
		if l < len(h.ary) && h.ary[l] < h.ary[minInd] {
			minInd = l
		}
		if r < len(h.ary) && h.ary[r] < h.ary[minInd] {
			minInd = r
		}
		if minInd != i {
			datastructures.Swap(h.ary, minInd, i)
			i = minInd
		} else {
			break
		}
	}
	img = append(img, arrayToColors(h.ary, max))
	return nil, min, img
}

func HeapSort(src []int) [][]color.RGBA {
	max := maxInt(src)
	var img [][]color.RGBA
	var heap InstrumentedHeap
	for i, ele := range src {
		imgSlice := heap.Insert(ele, max)

		var unsortedColors []color.RGBA
		for _, ele := range src[i+1:] {
			unsortedColors = append(unsortedColors, viridis(ele, max))
		}
		for i, _ := range imgSlice {
			imgSlice[i] = append(unsortedColors, imgSlice[i]...)
		}
		img = append(img, imgSlice...)
	}
	var sorted []int
	err, ele, imgSlice := heap.Pop(max)
	sorted = append(sorted, ele)
	for err == nil {
		var sortedColors []color.RGBA
		for _, ele := range sorted {
			sortedColors = append(sortedColors, viridis(ele, max))
		}
		for i, _ := range imgSlice {
			imgSlice[i] = append(sortedColors, imgSlice[i]...)
		}
		img = append(img, imgSlice...)

		err, ele, imgSlice = heap.Pop(max)
		sorted = append(sorted, ele)
	}
	copy(src, sorted)

	return img
}
