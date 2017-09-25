package sorting

import "image/color"

func BitonicSort(logn uint, a []int) [][]color.RGBA {
	if len(a) != 1 << logn {
		panic(1)
	}

	var img [][]color.RGBA
	max := maxInt(a)
	for i := 0; i < int(logn); i++ {
		for j := 0; j <= i; j++ {
			img = append(img, kernel(a, i, j, max)...)
		}
	}
	return img
}

func kernel(a []int, p int, q int, max int) [][]color.RGBA {
	d := 1 << uint(p-q)
	var img [][]color.RGBA

	for i := 0; i < len(a); i++ {
		up := ((i >> uint(p)) & 2) == 0
		if (i & d) == 0 && (a[i] > a[i | d]) == up {
			t := a[i]
			a[i] = a[i | d]
			a[i | d] = t
		}
		img = append(img, arrayToColors(a, max))
	}
	return img
}
