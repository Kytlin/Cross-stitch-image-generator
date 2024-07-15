package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/colormath"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
)

func GenerateGrid(img image.Image, threadColors []common.ThreadColor) [][]string {
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

func GenerateColorGrid(img image.Image, threadColors []common.ThreadColor, getNearestColor bool) [][]common.ThreadColor {
	bounds := img.Bounds()
	grid := make([][]common.ThreadColor, bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]common.ThreadColor, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y).(color.RGBA)

			if getNearestColor {
				row[x] = colormath.NearestColor(color, threadColors)
			} else {
				row[x] = common.ThreadColor{ // symbol necessary here?
					Color: color,
				}
			}
		}
		grid[y] = row
	}

	return grid
}

func findThreadColorName(c color.RGBA, threadColors []common.ThreadColor) string {
	for _, tc := range threadColors {
		if tc.Color == c {
			return tc.Name
		}
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// func findThreadColorObj(c color.RGBA, threadColors []common.ThreadColor) common.ThreadColor {
// 	for _, tc := range threadColors {
// 		if tc.Color == c {
// 			return tc
// 		}
// 	}
// 	return colormath.NearestColor(c, threadColors)
// }

func GetPartialPalette(img image.Image, threadColors []common.ThreadColor, k int) []common.ThreadColor {
	colorCounts := make(map[common.ThreadColor]int)
	bounds := img.Bounds()

	// Count the occurrence of each color in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			threadColor := colormath.NearestColor(c, threadColors)
			colorCounts[threadColor]++
		}
	}

	// Create a slice of thread colors and sort by count
	type colorCountPair struct {
		Color common.ThreadColor
		Count int
	}

	var sortedColors []colorCountPair
	for color, count := range colorCounts {
		sortedColors = append(sortedColors, colorCountPair{Color: color, Count: count})
	}
	sort.Slice(sortedColors, func(i, j int) bool {
		return sortedColors[i].Count > sortedColors[j].Count
	})

	// Select up to k colors
	var selectedColors []common.ThreadColor
	for i := 0; i < len(sortedColors) && i < k; i++ {
		selectedColors = append(selectedColors, sortedColors[i].Color)
	}

	return selectedColors
}
