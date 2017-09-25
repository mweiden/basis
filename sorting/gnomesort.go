package sorting

import (
	"image/color"
)

func GnomeSort(src []int) [][]color.RGBA {
	max := maxInt(src)
	var img [][]color.RGBA
	ind := 0
	for ind < len(src) {
		img = append(img, arrayToColors(src, max))
		if ind == 0 || src[ind] >= src[ind-1] {
			ind++
		} else {
			tmp := src[ind]
			src[ind] = src[ind-1]
			src[ind-1] = tmp
			ind--
		}
	}
	img = append(img, arrayToColors(src, max))
	return img
}

func arrayToColors(ary []int, max int) []color.RGBA {
	var colors []color.RGBA
	for _, ele := range ary {
		colors = append(colors, viridis(ele, max))
	}
	return colors
}
