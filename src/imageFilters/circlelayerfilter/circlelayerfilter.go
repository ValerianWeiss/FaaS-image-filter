package main

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"fmt"
	"image"
	"io/ioutil"
	"log"
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
	bimg, cimg, ftype := getImgs(req)
	r := calcRadius(cimg)

	newImg := addLayerCircle(bimg, cimg, r)

	return utils.CreateResJSON(newImg, ftype)
}

// addLayerCircle creates a circle out of the cImg and lays it over the base image
func addLayerCircle(baseImg image.Image, cImg image.Image, r float64) *image.RGBA {
	newImg := utils.Copy(baseImg)
	cWidth, cHeight := utils.GetImgSize(cImg)
	width, height := utils.GetImgSize(baseImg)
	inPosX := (width - cWidth) / 2
	inPosY := (height - cHeight) / 2
	centerX := cWidth / 2
	centerY := cHeight / 2
	for x := 0; x < cWidth; x++ {
		for y := 0; y < cHeight; y++ {
			if utils.Dist(float64(x), float64(y), float64(centerX), float64(centerY)) < r {
				c := cImg.At(x, y)
				newImg.Set(inPosX+x, inPosY+y, c)
			}
		}
	}
	return newImg
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

func calcRadius(cimg image.Image) float64 {
	var r float64
	cwidth, cheight := utils.GetImgSize(cimg)

	if cwidth >= cheight {
		r = float64(cheight) / 2
	} else {
		r = float64(cwidth) / 2
	}
	return r
}
