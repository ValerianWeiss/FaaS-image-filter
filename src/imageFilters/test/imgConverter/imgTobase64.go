package main

import (
	"encoding/base64"
	"log"
	"os"
)

func main() {
	filePath := os.Args[0]
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	fInfo, _ := file.Stat()
	fsize := fInfo.Size()
	buffer := make([]byte, fsize)
	file.Read(buffer)

	imgBase64Str := base64.StdEncoding.EncodeToString(buffer)
	strFile, _ := os.Create("./imgBase64Str.txt")
	strFile.Write([]byte(imgBase64Str))
}

