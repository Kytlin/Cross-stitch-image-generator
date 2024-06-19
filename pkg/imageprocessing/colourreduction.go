package imageprocessing

import (
	"bufio"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ThreadColour struct {
	ID     int
	Name   string
	Colour color.RGBA
}

func ColourAtoi(s string) uint8 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("error converting string to int: %s", err)
		return 0
	}
	return uint8(i)
}

func LoadThreadColours(filePath string) ([455]ThreadColour, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var threadImg [455]ThreadColour
	threadIdx := 0
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

		threadColour := ThreadColour{
			ID:     id,
			Name:   name,
			Colour: color.RGBA{R: ColourAtoi(parts[lineIdx-1]), G: ColourAtoi(parts[lineIdx]), B: ColourAtoi(parts[lineIdx+1])},
		}

		threadImg[threadIdx] = threadColour
		threadIdx += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	return threadImg, err
}

func ReduceColors(img image.Image, palette [455]ThreadColour) image.Image {
	bounds := img.Bounds()
	reducedImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 1 {
			originalColour := img.At(x, y)
			nearestColour := findNearestColour(originalColour, palette)
			reducedImg.Set(x, y, nearestColour)
		}
	}

	return reducedImg
}

func findNearestColour(originalColour color.Color, palette [455]ThreadColour) color.Color {
	rImg, gImg, bImg, _ := originalColour.RGBA()
	rImg, gImg, bImg = rImg>>8, gImg>>8, bImg>>8
	minDistance := math.MaxFloat64
	var nearestColor color.Color

	for _, threadColour := range palette {
		rPal := uint8(threadColour.Colour.R)
		gPal := uint8(threadColour.Colour.G)
		bPal := uint8(threadColour.Colour.B)

		distance := colorDistance(rImg, gImg, bImg, rPal, gPal, bPal)
		if distance < minDistance {
			minDistance = distance
			nearestColor = color.RGBA{
				R: rPal,
				G: gPal,
				B: bPal,
			}
		}
	}

	return nearestColor
}

func colorDistance(rImg, gImg, bImg uint32, rPal, gPal, bPal uint8) float64 {
	return math.Sqrt(
		float64((rImg-uint32(rPal))*(rImg-uint32(rPal)) +
			(gImg-uint32(gPal))*(gImg-uint32(gPal)) +
			(bImg-uint32(bPal))*(bImg-uint32(bPal))))
}
