package main

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
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

func decodeImg(res *http.Response, err error) (image.Image, string) {

	if err != nil {
		panic(err)
	}

	imgBase64str := getImgStrFromBody(res)
	return utils.DecodeBase64Img(imgBase64str)
}

func getImgStrFromBody(res *http.Response) string {
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	resJSONMap := make(map[string]interface{})
	json.Unmarshal(body, &resJSONMap)
	return resJSONMap["image"].(string)
}

func getImgs(req []byte) (image.Image, image.Image, string) {
	var wg sync.WaitGroup
	wg.Add(2)

	var bimg image.Image
	var cimg image.Image
	var ftype string

	go func() {
		reader := bytes.NewReader(req)
		res, err := http.Post("http://127.0.0.1:8080/function/blurfilter", "application/json", reader)
		bimg, ftype = decodeImg(res, err)
		wg.Done()
	}()

	go func() {
		reqJSONMap := make(map[string]interface{})
		json.Unmarshal(req, &reqJSONMap)
		reqJSONMap["scaling"] = 0.9
		scaleReq, _ := json.Marshal(reqJSONMap)
		reader := bytes.NewReader(scaleReq)
		res, err := http.Post("http://127.0.0.1:8080/function/scalefilter", "application/json", reader)
		cimg, _ = decodeImg(res, err)
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
