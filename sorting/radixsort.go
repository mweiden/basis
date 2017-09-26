package sorting

import (
	"image/color"
)

func RadixSort(src []int) [][]color.RGBA {
	max := maxInt(src)
	var img [][]color.RGBA
	sortingStack := [][]int{src}

	img = append(img, sortingStackToColors(sortingStack, max))

	for k := 31; k >= 0; k-- {
		var tmp [][]int
		for _, ary := range sortingStack {
			zeros, ones := split(ary, uint(k))
			if len(zeros) > 0 {
				tmp = append(tmp, zeros)
			}
			if len(ones) > 0 {
				tmp = append(tmp, ones)
			}
		}
		sortingStack = tmp
		img = append(img, sortingStackToColors(sortingStack, max))
	}
	ind := 0
	for _, ary := range sortingStack {
		for _, ele := range ary {
			src[ind] = ele
			ind++
		}
	}
	return img
}

func split(src []int, k uint) ([]int, []int) {
	if len(src) == 0 {
		return []int{}, []int{}
	}
	var zeros []int
	var ones []int
	for _, ele := range src {
		if selectBit(ele, k) == 0 {
			zeros = append(zeros, ele)
		} else {
			ones = append(ones, ele)
		}
	}
	return zeros, ones
}

func selectBit(i int, k uint) int {
	return (i & (1 << k)) >> k
}

func sortingStackToColors(sortingStack [][]int, n int) []color.RGBA {
	ints := flatten(sortingStack)
	var colors []color.RGBA
	for _, ele := range ints {
		colors = append(colors, viridis(ele, n))
	}
	return colors
}

func flatten(arys [][]int) []int {
	var flat []int
	for _, ary := range arys {
		for _, ele := range ary {
			flat = append(flat, ele)
		}
	}
	return flat
}

func maxInt(ints []int) int {
	max := ints[0]
	for _, ele := range ints[1:] {
		if ele > max {
			max = ele
		}
	}
	return max
}
