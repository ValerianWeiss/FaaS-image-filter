package imagefilter

import (
	"image"
	"image/color"

	"imagefilter.com/src/imageFunctions/utils"
)

// Handle a serverless request
func Handle(req []byte) image.RGBA {
	reqMap := utils.ParseJSON(req)
	imgBytes := reqMap["image"].([]byte)
	img := utils.ReadImgBytes(imgBytes)
	newImg := blackWhite(img)

	return *newImg
}

// BlackWhite is a image filter which creates a balck and white image of
// a given image and returns it
func blackWhite(img image.Image) *image.RGBA {
	width, height := utils.GetImgSize(img)
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			gray := uint8((r + g + b) / 3 / 257)
			nr := gray
			ng := gray
			nb := gray
			na := uint8(a / 257)

			c := color.RGBA{nr, ng, nb, na}
			newImg.Set(x, y, c)
		}
	}
	return newImg
}
