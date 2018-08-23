package main

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	ftype := os.Args[1][strings.Index(os.Args[1], "."):]
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
	imgBase64Str := "data:image/" + ftype + ";base64," + base64.StdEncoding.EncodeToString(buffer)

	reqMap := map[string]string{"image": imgBase64Str}
	req, _ := json.Marshal(reqMap)
	reqReader := bytes.NewReader(req)

	res, err := http.Post("127.0.0.1:8080/", "application/json", reqReader)

	if err != nil {
		panic(err)
	}

	var resBuffer []byte
	res.Body.Read(resBuffer)

	jsonMap := make(map[string]interface{})
	json.Unmarshal(resBuffer, &jsonMap)

	newImgBase64str := jsonMap["image"].(string)

	img, ftype := utils.DecodeBase64Img(newImgBase64str)
	enImg := utils.EncodeBase64Img(img, ftype)

	nf, _ := os.Create("./out." + ftype)
	nf.Write([]byte(enImg))
	defer nf.Close()
}
