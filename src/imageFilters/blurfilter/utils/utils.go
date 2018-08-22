package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"strings"
)

// ParseJSON parses a json object to a map
func ParseJSON(jsonObj []byte) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	e := json.Unmarshal(jsonObj, &jsonMap)

	test := jsonMap["image"].(string)
	panic(test)
	//panic on error
	if e != nil {
		panic(e)
	} else {
		return jsonMap
	}
}

// ReadImgBytes reads and decodes an image from a byte array
func ReadImgBytes(content []byte) image.Image {
	reader := bytes.NewReader(content)
	registerFormats()
	img, _, _ := image.Decode(reader)
	return img
}

// GetImgSize returns the size of an image
func GetImgSize(img image.Image) (width int, height int) {
	size := img.Bounds()
	return size.Max.X, size.Max.Y
}

// Dist calculates the distance between two points
func Dist(x1, y1, x2, y2 float64) float64 {
	len1 := math.Sqrt(math.Pow(x1, 2.0) + math.Pow(y1, 2.0))
	len2 := math.Sqrt(math.Pow(x2, 2.0) + math.Pow(y2, 2.0))

	var vx float64
	var vy float64

	if len1 > len2 {
		vx = x1 - x2
		vy = y1 - y2
	} else if len1 < len2 {
		vx = x2 - x1
		vy = y2 - y1
	} else {
		return 0
	}
	return math.Sqrt(math.Pow(vx, 2.0) + math.Pow(vy, 2.0))
}

// Registers the jpg and png decoder
func registerFormats() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

// GetAvgColor calculates the average color in of pixels which are left and right
// of the given position
func GetAvgColor(x int, y int, img image.Image, xrange int) color.RGBA {

	width, _ := GetImgSize(img)
	var iterationCount uint32
	values := []uint32{0, 0, 0, 0}

	for ix := x - xrange/2; ix < x+xrange/2; ix++ {
		if ix >= 0 && ix < width {
			r, g, b, a := img.At(ix, y).RGBA()
			values[0] += r / 257
			values[1] += g / 257
			values[2] += b / 257
			values[3] += a / 257
			iterationCount++
		}
	}

	for i := 0; i < len(values); i++ {
		values[i] = values[i] / iterationCount
	}
	return color.RGBA{uint8(values[0]), uint8(values[1]), uint8(values[2]), uint8(values[3])}
}

// DecodeBase64Img decodes a image which is encoded as base64 string
func DecodeBase64Img(base64str string) image.Image {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64str))
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	return img
}
