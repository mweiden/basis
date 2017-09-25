package main

import (
	"image"
	"image/png"
	"math/rand"
	"os"
)

func main() {
	//logn := uint(9)
	//height := 1 << logn
	height := 1500
	colorInds := rand.Perm(height)
	imgData := sorting.HeapSort(colorInds)

	iters := len(imgData)

	heightMult := 6
	widthMult := 1

	// Create an image
	img := image.NewRGBA(image.Rect(0, 0, iters*widthMult, height*heightMult))

	for i := 0; i < iters; i++ {
		for j := 0; j < height; j++ {
			for ii := 0; ii < widthMult; ii++ {
				for jj := 0; jj < heightMult; jj++ {
					img.Set(i*widthMult+ii, j*heightMult+jj, imgData[i][j])
				}
			}
		}
	}

	// Save to out.png
	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)

	last := -1
	for _, ele := range colorInds {
		if ele <= last {
			os.Exit(1)
		}
		last = ele
	}
}
