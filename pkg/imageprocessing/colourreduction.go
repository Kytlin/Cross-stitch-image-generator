package imageprocessing

import (
	"bufio"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/colourmath"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
)

func ColourAtoi(s string) uint8 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("error converting string to int: %s", err)
		return 0
	}
	return uint8(i)
}

func LoadThreadColours(filePath string) ([]common.ThreadColour, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var threadImg []common.ThreadColour
	for scanner.Scan() {
		line := scanner.Text()

		// Find the last space before the integer
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '\t'
		})

		var name string
		var id int
		lineIdx := len(parts) - 3

		id, err = strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		name = strings.Join(parts[1:lineIdx], " ")

		ThreadColour := common.ThreadColour{
			ID:     id,
			Name:   name,
			Colour: color.RGBA{R: ColourAtoi(parts[lineIdx-1]), G: ColourAtoi(parts[lineIdx]), B: ColourAtoi(parts[lineIdx+1])},
		}

		threadImg = append(threadImg, ThreadColour)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	return threadImg, err
}

func ReduceColors(img image.Image, palette []common.ThreadColour) image.Image {
	bounds := img.Bounds()
	reducedImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 1 {
			originalColour := img.At(x, y)
			nearestColour := colourmath.NearestColour(originalColour, palette)
			reducedImg.Set(x, y, nearestColour.Colour)
		}
	}

	return reducedImg
}
