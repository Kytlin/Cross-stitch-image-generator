package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
)

func GenerateGrid(img image.Image, threadColors [455]ThreadColour) [][]string {
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

func findThreadColorName(c color.RGBA, threadColors [455]ThreadColour) string {
	for _, tc := range threadColors {
		if tc.Colour == c {
			return tc.Name
		}
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
