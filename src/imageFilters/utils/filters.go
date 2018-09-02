package utils

import (
	"bytes"
	"encoding/json"
	"image"
	"net/http"
)

// BlurImg expects a JSON byte array with a single field which is named "image"
// which contains the base64 encoded image. The function will add the given "blurscale"
// parameter to the request and returns the result image and the file type
func BlurImg(req []byte, blurscale int) (image.Image, string) {
	reqJSONMap := make(map[string]interface{})
	json.Unmarshal(req, &reqJSONMap)
	reqJSONMap["blurscale"] = blurscale
	blurReq, _ := json.Marshal(reqJSONMap)
	reader := bytes.NewReader(blurReq)
	res, err := http.Post("http://gateway:8080/function/blurfilter", "application/json", reader)
	return DecodeImg(res, err)
}

// ScaleImg expects a JSON byte array with a single field which is named "image"
// which contains the base64 encoded image. The function will add the given "scaling"
// parameter to the request and returns the result image and the file type
func ScaleImg(req []byte, scaling float64) (image.Image, string) {
	reqJSONMap := make(map[string]interface{})
	json.Unmarshal(req, &reqJSONMap)
	reqJSONMap["scaling"] = scaling
	scaleReq, _ := json.Marshal(reqJSONMap)
	reader := bytes.NewReader(scaleReq)
	res, err := http.Post("http://gateway:8080/function/scalefilter", "application/json", reader)
	return DecodeImg(res, err)
}
