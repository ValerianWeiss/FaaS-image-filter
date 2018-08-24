package main

import (
	"blurfilter/utils"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
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

// Handle a serverless request
func handle(req []byte) string {
	jsonMap := utils.ParseJSON(req)
	imgBase64str := jsonMap["image"].(string)
	blurscale := int(jsonMap["blurscale"].(float64))

	img, imgType := utils.DecodeBase64Img(imgBase64str)
	newImg := blur(img, blurscale)
	newImgBase64str := utils.EncodeBase64Img(newImg, imgType)

	resMap := map[string]string{"image": newImgBase64str}
	res, _ := json.Marshal(resMap)

	return string(res)
}

// blur blures an image
func blur(img image.Image, gridSize int) *image.RGBA {
	width, height := utils.GetImgSize(img)
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	gridLen := math.Floor(math.Sqrt(float64(gridSize)))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := utils.GetAvgColor(x, y, img, int(gridLen))

			for ix := x - int(gridLen/2); ix < x+int(gridLen/2); ix++ {
				if ix >= 0 && ix < width {
					newImg.Set(ix, y, c)
				}
			}
		}
	}
	return newImg
}
