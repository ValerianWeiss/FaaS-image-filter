package utils

import (
	"image"
	"math"
)

// ScaleDown scales a image down to a given size
func ScaleDown(img image.Image, width int, height int) *image.RGBA {
	owidth, oheight := GetImgSize(img)
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	setpWidthX := owidth / width
	setpWidthY := oheight / height

	for x, ox := 0, 0; x < width; x, ox = x+1, ox+setpWidthX {
		for y, oy := 0, 0; y < height; y, oy = y+1, oy+setpWidthY {
			c := img.At(ox, oy)
			newImg.Set(x, y, c)
		}
	}
	return newImg
}

// AddLayerTriangle creates a triangle out of the trImg and lays it over the base image
func AddLayerTriangle(baseImg image.Image, trImg image.Image, edgeLen int) *image.RGBA {
	newImg := Copy(baseImg)
	trWidth, trHeight := GetImgSize(trImg)
	width, height := GetImgSize(baseImg)
	inPosX := (width - trWidth) / 2
	inPosY := (height - trHeight) / 2
	centerX := trWidth / 2
	centerY := trHeight / 2
	h := int(math.Sqrt(math.Pow(float64(edgeLen), 2.0) + math.Pow(float64(edgeLen/2), 2.0)))
	p1 := []int{centerX - edgeLen/2, centerY + h/2}
	p2 := []int{centerX, centerY - h/2}
	p3 := []int{centerX + edgeLen/2, centerY + h/2}

	f1 := func(x int) int {
		return ((p1[1]-p2[1])/(p1[0]-p2[0]))*x + (p2[1]*p1[0]-p1[1]*p2[0])/(p1[0]-p2[0])
	}

	f2 := func(x int) int {
		return ((p3[1]-p2[1])/(p3[0]-p2[0]))*x + (p2[1]*p3[0]-p3[1]*p2[0])/(p3[0]-p2[0])
	}

	for x := 0; x < trWidth; x++ {
		for y := 0; y < trHeight; y++ {
			if y > f1(x) && y > f2(x) && y < centerY+h/2 {
				c := trImg.At(x, y)
				newImg.Set(inPosX+x, inPosY+y, c)
			}
		}
	}
	return newImg
}

// AddLayerCircle creates a circle out of the cImg and lays it over the base image
func AddLayerCircle(baseImg image.Image, cImg image.Image, r float64) *image.RGBA {
	newImg := Copy(baseImg)
	cWidth, cHeight := GetImgSize(cImg)
	width, height := GetImgSize(baseImg)
	inPosX := (width - cWidth) / 2
	inPosY := (height - cHeight) / 2
	centerX := cWidth / 2
	centerY := cHeight / 2

	for x := 0; x < cWidth; x++ {
		for y := 0; y < cHeight; y++ {
			if Dist(float64(x), float64(y), float64(centerX), float64(centerY)) < r {
				c := cImg.At(x, y)
				newImg.Set(inPosX+x, inPosY+y, c)
			}
		}
	}
	return newImg
}

// Copy copies a image
func Copy(img image.Image) *image.RGBA {
	width, height := GetImgSize(img)
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := img.At(x, y)
			newImg.Set(x, y, c)
		}
	}
	return newImg
}
