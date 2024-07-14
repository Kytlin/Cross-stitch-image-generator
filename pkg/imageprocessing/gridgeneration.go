package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/colourmath"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
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

func GenerateColourGrid(img image.Image, threadColors []common.ThreadColour) [][]common.ThreadColour {
	bounds := img.Bounds()
	grid := make([][]common.ThreadColour, bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]common.ThreadColour, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y).(color.RGBA)

			row[x] = colourmath.NearestColour(color, threadColors)
			fmt.Println(row[x].Colour)
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

func GetPartialPalette(img image.Image, threadColours []common.ThreadColour, k int) []common.ThreadColour {
	colorCounts := make(map[common.ThreadColour]int)
	bounds := img.Bounds()

	// Count the occurrence of each color in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			threadColor := colourmath.NearestColour(c, threadColours)
			colorCounts[threadColor]++
		}
	}

	// Create a slice of thread colors and sort by count
	type colorCountPair struct {
		Color common.ThreadColour
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
	var selectedColors []common.ThreadColour
	for i := 0; i < len(sortedColors) && i < k; i++ {
		selectedColors = append(selectedColors, sortedColors[i].Color)
	}

	return selectedColors
}
