package main

import (
	"FaaS-image-filter/src/imageFilters/test/utils"
	"os"
)

func main() {
	fpath := os.Args[1]
	reqMap := map[string]interface{}{}
	url := "http://127.0.0.1:8080/function/circlelayerfilter"
	utils.ExecFunc(fpath, url, reqMap)
}
