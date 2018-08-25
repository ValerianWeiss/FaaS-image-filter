package utils

import (
	"FaaS-image-filter/src/imageFilters/utils"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ExecFunc executes a given imagefilter function and stores the result image
// as result.png for example in the same directory. It also wirtes to request and
// response in .txt files
func ExecFunc(imgPath, url string, reqMap map[string]interface{}) {
	f, err := os.Open(imgPath)
	ftype := imgPath[strings.Index(imgPath, ".")+1:]
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
	reqMap["image"] = imgBase64Str

	req, _ := json.Marshal(reqMap)
	reqReader := bytes.NewReader(req)

	tf, _ := os.Create("./request.txt")
	tf.Write(req)
	defer tf.Close()

	res, err := http.Post(url, "application/json", reqReader)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	jsonMap := make(map[string]interface{})
	e := json.Unmarshal(body, &jsonMap)

	of, _ := os.Create("./response.txt")
	of.Write(body)

	if e != nil {
		fmt.Println("Could not parse response")
		panic(e)
	}

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
