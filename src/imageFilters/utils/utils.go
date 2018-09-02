package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
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

// Copy copies a image to a new RGBA image
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

// CreateResJSON creates a JSON where the image is getting base64 encoded
// and stored unter the "image" fied of the JSON
func CreateResJSON(resImg image.Image, ftype string) string {
	newImgBase64str := EncodeBase64Img(resImg, ftype)
	resMap := map[string]string{"image": newImgBase64str}
	res, _ := json.Marshal(resMap)
	return string(res)
}

// DecodeImg gets the string property of "image" of the http response,
// which has to be a base64 encoded image and decodes it
func DecodeImg(res *http.Response, err error) (image.Image, string) {

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	resJSONMap := make(map[string]interface{})
	json.Unmarshal(body, &resJSONMap)

	imgBase64str := resJSONMap["image"].(string)
	return DecodeBase64Img(imgBase64str)
}
