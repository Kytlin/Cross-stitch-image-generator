package colourmath

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

func NearestColour(originalColour color.Color, palette []common.ThreadColour) common.ThreadColour {
	rImg, gImg, bImg, _ := originalColour.RGBA()
	rImg, gImg, bImg = rImg>>8, gImg>>8, bImg>>8
	minDistance := math.MaxFloat64
	var nearestColour common.ThreadColour

	for _, threadColour := range palette {
		rPal := uint8(threadColour.Colour.R)
		gPal := uint8(threadColour.Colour.G)
		bPal := uint8(threadColour.Colour.B)

		distance := colorDistance(rImg, gImg, bImg, rPal, gPal, bPal)
		if distance < minDistance {
			minDistance = distance
			nearestColour = threadColour
		}
	}

	return nearestColour
}
