package imagefilter

import (
	"image"

	filters "imagefilter.com/src/filters"
	utils "imagefilter.com/src/filterutils"
)

func main() {
	intputValid, iimgPath, iimgType, oimgPath, oimgType := utils.CheckPathArguments()

	if intputValid {
		img := utils.DecodeImg(iimgPath, iimgType)
		width, height := utils.GetImgSize(img)
		var newImg *image.RGBA

		newImg = filters.Blur(img, 500)
		triangleImg := filters.ScaleDown(img, int(0.9*float32(width)), int(0.9*float32(height)))
		newImg = filters.AddLayerTriangle(newImg, triangleImg, 400)

		utils.CreateImgFile(oimgPath, newImg, oimgType)
	}
}
