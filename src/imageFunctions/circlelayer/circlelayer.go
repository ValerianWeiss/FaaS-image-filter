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
		var r float64

		newImg = filters.Blur(img, 500)
		circleImg := filters.ScaleDown(img, int(0.9*float32(width)), int(0.9*float32(height)))

		cwidth, cheight := utils.GetImgSize(circleImg)

		if cwidth >= cheight {
			r = float64(cheight) / 2
		} else {
			r = float64(cwidth) / 2
		}

		newImg = filters.AddLayerCircle(newImg, circleImg, r)
		utils.CreateImgFile(oimgPath, newImg, oimgType)
	}
}
