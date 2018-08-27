package main

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
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
	jsonMap := utils.ParseJSON(req)
	imgBase64str := jsonMap["image"].(string)
	scaling := jsonMap["image"].(float64)
	img, ftype := utils.DecodeBase64Img(imgBase64str)

	scaledImg := scale(img, scaling)
	newImgBase64str := utils.EncodeBase64Img(scaledImg, ftype)
	resMap := map[string]string{"image": newImgBase64str}
	res, _ := json.Marshal(resMap)
	return string(res)
}

// scale scales a image by a given factor
func scale(img image.Image, scaling float64) *image.RGBA {
	owidth, oheight := utils.GetImgSize(img)
	width := owidth * int(scaling)
	height := oheight * int(scaling)
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	setpWidthX := float64(owidth) / float64(width)
	setpWidthY := float64(oheight) / float64(height)

	for x, ox := 0, 0.0; x < width; x, ox = x+1, ox+setpWidthX {
		for y, oy := 0, 0.0; y < height; y, oy = y+1, oy+setpWidthY {
			c := img.At(int(ox), int(oy))
			newImg.Set(x, y, c)
		}
	}
	return newImg
}
