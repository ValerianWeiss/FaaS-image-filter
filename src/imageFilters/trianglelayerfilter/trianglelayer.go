package main

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		log.Fatalf("Unable to read standard input: %s", err.Error())
	}

	output := handle(input)
	fmt.Println(output)
}

func handle(req []byte) string {
	bImg, trImg, ftype := getImgs(req)
	newImg := addLayerTriangle(bImg, trImg)

	return utils.CreateResJSON(newImg, ftype)
}

// addLayerTriangle creates a triangle out of the trImg and lays it over the base image
func addLayerTriangle(baseImg image.Image, trImg image.Image) *image.RGBA {
	trWidth, trHeight := utils.GetImgSize(trImg)
	width, height := utils.GetImgSize(baseImg)
	newImg := utils.Copy(baseImg)
	inPosX := (width - trWidth) / 2
	inPosY := (height - trHeight) / 2
	centerX := trWidth / 2
	centerY := trHeight / 2
	edgeLen := calcEdgeLen(trWidth, trHeight)

	h := int(math.Sqrt(math.Pow(float64(edgeLen), 2) - math.Pow(float64(edgeLen/2), 2)))
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
			if y > f1(x) && y > f2(x) {
				c := trImg.At(x, y)
				newImg.Set(inPosX+x, inPosY+y, c)
			}
		}
	}
	return newImg
}

func calcEdgeLen(width, height int) int {
	if width > height {
		angle := 2 * math.Pi / 360 * 30
		return int(float64(height) / math.Cos(angle))
	}
	return width
}

func getImgs(req []byte) (image.Image, image.Image, string) {
	var wg sync.WaitGroup
	wg.Add(2)

	var bimg image.Image
	var cimg image.Image
	var ftype string

	go func() {
		bimg, ftype = utils.BlurImg(req, 500)
		wg.Done()
	}()

	go func() {
		cimg, _ = utils.ScaleImg(req, 0.9)
		wg.Done()
	}()

	wg.Wait()
	return bimg, cimg, ftype
}
