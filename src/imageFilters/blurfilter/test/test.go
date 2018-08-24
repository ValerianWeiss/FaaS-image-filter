package main

import (
	function "FaaS-image-filter/src/imageFilters/blurfilter"
	"FaaS-image-filter/src/imageFilters/utils"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	ftype := os.Args[1][strings.LastIndex(os.Args[1], ".")+1:]

	if ftype == "jpg" {
		ftype = "jpeg"
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create a new buffer base on file size
	fInfo, _ := f.Stat()
	fsize := fInfo.Size()
	buffer := make([]byte, fsize)

	// read file content into buffer
	fReader := bufio.NewReader(f)
	fReader.Read(buffer)

	// convert the buffer bytes to base64 string
	imgBase64Str := "data:image/" + ftype + ";base64," + base64.StdEncoding.EncodeToString(buffer)

	reqMap := map[string]interface{}{"image": imgBase64Str, "blurscale": 500}
	req, _ := json.Marshal(reqMap)
	//reqReader := bytes.NewReader(req)

	resBuffer := function.Handle(req)

	jsonMap := make(map[string]interface{})
	json.Unmarshal(resBuffer, &jsonMap)

	newImgBase64str := jsonMap["image"].(string)

	img, ftype := utils.DecodeBase64Img(newImgBase64str)

	nf, _ := os.Create("./out." + ftype)

	if ftype == "png" {
		png.Encode(nf, img)
	} else {
		jpeg.Encode(nf, img, nil)
	}
	defer nf.Close()
}
