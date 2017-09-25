package sorting

import "image/color"

func MergeSort(src []int) [][]color.RGBA {
	tmp := make([]int, len(src))
	copy(tmp, src)
	return bottomUpMergeSort(src, tmp)
}

func bottomUpMergeSort(A []int, B []int) [][]color.RGBA {
	n := len(A)
	var img [][]color.RGBA

	for width := 1; width < n; width *= 2 {
		for i := 0; i < n; i += 2 * width {
			imgSlice := bottomUpMerge(A, i, min(i+width, n), min(i+2*width, n), B)
			for _, colors := range imgSlice {
				img = append(img, colors)
			}
		}
		tmp := A
		A = B
		B = tmp
		//copy(A,B)
	}
	copy(B, A)
	return img
}

func bottomUpMerge(A []int, iLeft int, iRight int, iEnd int, B []int) [][]color.RGBA {
	left := iLeft
	right := iRight
	var imgSlice [][]color.RGBA

	for k := iLeft; k < iEnd; k++ {
		if left < iRight && (right >= iEnd || A[left] <= A[right]) {
			B[k] = A[left]
			left += 1
		} else {
			B[k] = A[right]
			right += 1
		}
		imgSlice = append(imgSlice, getColors(B, k))
	}
	imgSlice = append(imgSlice, getColors(B, iEnd))

	return imgSlice
}

func getColors(ary []int, k int) []color.RGBA {
	n := len(ary)
	colors := make([]color.RGBA, n)
	for i, ele := range ary {
		colors[i] = viridis(ele, n)
		if i > k {
			colors[i].R = uint8(float64(colors[i].R) * 0.5)
			colors[i].G = uint8(float64(colors[i].G) * 0.5)
			colors[i].B = uint8(float64(colors[i].B) * 0.5)
		}
	}
	return colors
}

func min(i int, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}
