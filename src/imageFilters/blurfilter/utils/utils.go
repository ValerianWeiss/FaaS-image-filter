package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
)

// ParseJSON parses a json object to a map
func ParseJSON(jsonObj []byte) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	e := json.Unmarshal(jsonObj, &jsonMap)

	//panic on error
	if e != nil {
		panic(e)
	} else {
		return jsonMap
	}
}

// GetImgSize returns the size of an image
func GetImgSize(img image.Image) (width int, height int) {
	size := img.Bounds()
	return size.Max.X, size.Max.Y
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
func DecodeBase64Img(base64str string) (image.Image, string) {

	var unbased []byte
	var err error
	var img image.Image
	soi := strings.Index(base64str, ",")
	types := strings.Index(base64str, "/") + 1
	typee := strings.Index(base64str, ";")
	filetype := base64str[types:typee]

	unbased, err = base64.StdEncoding.DecodeString(base64str[soi+1:])

	if err != nil {
		panic("Cannot decode b64")
	}

	reader := bytes.NewReader(unbased)

	if filetype == "jpeg" {
		img, err = jpeg.Decode(reader)
	} else if filetype == "png" {
		img, err = png.Decode(reader)
	} else {
		panic("Filetype " + filetype + " not supported\nSupported filetypes: jpeg, png")
	}

	if err != nil {
		panic(err)
	}

	return img, filetype
}

// EncodeBase64Img ecodes a  image to base64
func EncodeBase64Img(img image.Image, filetype string) string {
	buffer := new(bytes.Buffer)
	err := jpeg.Encode(buffer, img, nil)

	if err != nil {
		panic(err)
	}

	// convert the buffer bytes to base64 string
	imgBase64Str := "data:image/" + filetype + ";base64," + base64.StdEncoding.EncodeToString(buffer.Bytes())
	return imgBase64Str
}
