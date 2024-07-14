package imageprocessing

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/colourmath"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
	"github.com/ericpauley/go-quantize/quantize"
)

func GenerateGrid(img image.Image, threadColors []common.ThreadColour) [][]string {
	bounds := img.Bounds()
	grid := make([][]string, bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]string, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y).(color.RGBA)
			row[x] = findThreadColorName(color, threadColors)
		}
		grid[y] = row
	}

	return grid
}

func GenerateColourGrid(img image.Image, threadColors []common.ThreadColour, getNearestColour bool) [][]common.ThreadColour {
	bounds := img.Bounds()
	grid := make([][]common.ThreadColour, bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]common.ThreadColour, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y).(color.RGBA)

			if getNearestColour {
				row[x] = colourmath.NearestColour(color, threadColors)
				fmt.Println(row[x].Colour)
			} else {
				row[x] = common.ThreadColour{
					Colour: color,
					Symbol: "a",
				}
			}
		}
		grid[y] = row
	}

	return grid
}

func findThreadColorName(c color.RGBA, threadColors []common.ThreadColour) string {
	for _, tc := range threadColors {
		if tc.Colour == c {
			return tc.Name
		}
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// func findThreadColorObj(c color.RGBA, threadColors []common.ThreadColour) common.ThreadColour {
// 	for _, tc := range threadColors {
// 		if tc.Colour == c {
// 			return tc
// 		}
// 	}
// 	return colourmath.NearestColour(c, threadColors)
// }

// GetPartialPalette generates a palette of k colors from the image using available thread colors
func GetPartialPalette(img image.Image, threadColours []common.ThreadColour, k int) []common.ThreadColour {
	quantizer := quantize.MedianCutQuantizer{}
	palette := quantizer.Quantize(make([]color.Color, 0, k), img)

	// Convert quantized palette to the nearest thread colors
	threadPalette := make([]common.ThreadColour, len(palette))
	for i, c := range palette {
		threadPalette[i] = colourmath.NearestColour(c, threadColours)
	}

	return threadPalette
}
