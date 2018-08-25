package imagefilter

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"FaaS-image-filter/src/imageFilters/utils"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		log.Fatalf("Unable to read standard input: %s", err.Error())
	}

	output := handle(input)
	fmt.Println(output)
}

// handle a serverless request
func handle(req []byte) string {
	jsonMap := utils.ParseJSON(req)
	imgBase64str := jsonMap["image"].(string)

	img, imgType := utils.DecodeBase64Img(imgBase64str)
	newImg := blackWhite(img)
	newImgBase64str := utils.EncodeBase64Img(newImg, imgType)

	resMap := map[string]string{"image": newImgBase64str}
	res, _ := json.Marshal(resMap)

	return string(res)
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
