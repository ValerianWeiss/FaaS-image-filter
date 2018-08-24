package main

import (
	"blurfilter/utils"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"image/png"
	"image/jpeg"
)

func main() {
	f, err := os.Open(os.Args[1])
	ftype := os.Args[1][strings.Index(os.Args[1], ".")+1:]
	defer f.Close()

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
	imgBase64Str := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buffer)

	reqMap := map[string]interface{}{"image":imgBase64Str, "blurscale":500}
	req, _ := json.Marshal(reqMap)
	reqReader := bytes.NewReader(req)

	tf, _ := os.Create("./request.txt")
	tf.Write(req)
	defer tf.Close()	

	res, _ := http.Post("http://127.0.1.1:8080/function/imageblur", "application/json", reqReader)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)


	jsonMap := make(map[string]interface{})
	e := json.Unmarshal(body, &jsonMap)

	if e != nil {
		fmt.Println("Could not parse response")
		panic(e)
	}

	of, _ := os.Create("./response.txt")
	of.Write(body)

	newImgBase64str := jsonMap["image"].(string)
	img, ftype := utils.DecodeBase64Img(newImgBase64str)

	nf, _ := os.Create("./result." + ftype)
	defer nf.Close()

	if ftype == "png" {
		png.Encode(nf, img)
	} else {
		jpeg.Encode(nf, img, nil)
	}
}
