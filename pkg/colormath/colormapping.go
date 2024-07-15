package colormath

import (
	"image/color"
	"math"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
)

func colorDistance(rImg, gImg, bImg uint32, rPal, gPal, bPal uint8) float64 {
	return math.Sqrt(
		float64((rImg-uint32(rPal))*(rImg-uint32(rPal)) +
			(gImg-uint32(gPal))*(gImg-uint32(gPal)) +
			(bImg-uint32(bPal))*(bImg-uint32(bPal))))
}

func NearestColor(originalColor color.Color, palette []common.ThreadColor) common.ThreadColor {
	rImg, gImg, bImg, _ := originalColor.RGBA()
	rImg, gImg, bImg = rImg>>8, gImg>>8, bImg>>8
	minDistance := math.MaxFloat64
	var nearestColor common.ThreadColor

	for _, threadColor := range palette {
		rPal := uint8(threadColor.Color.R)
		gPal := uint8(threadColor.Color.G)
		bPal := uint8(threadColor.Color.B)

		distance := colorDistance(rImg, gImg, bImg, rPal, gPal, bPal)
		if distance < minDistance {
			minDistance = distance
			nearestColor = threadColor
		}
	}

	return nearestColor
}
