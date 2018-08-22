package function

import (
	"image"
	"math"

	"handler/function/utils"
)

// Handle a serverless request
func Handle(req []byte) image.RGBA {
	jsonMap := utils.ParseJSON(req)
	imgBase64 := jsonMap["image"].(string)
	blurscale := int(jsonMap["blurscale"].(float64))

	img := utils.DecodeBase64Img(imgBase64)
	newImg := blur(img, blurscale)

	return *newImg
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