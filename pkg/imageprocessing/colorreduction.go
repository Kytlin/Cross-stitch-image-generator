package imageprocessing

import (
	"bufio"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/colormath"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
)

func createUnicodeCharMap(threadColors []common.ThreadColor) map[int]rune {
	unicodeMap := make(map[int]rune)
	mapIdx := 0
	threadColorsLength := len(threadColors)

	ranges := [][3]int{
		{0x2190, 0x21FF}, // Arrows: U+2190 to U+21FF
		{0x2200, 0x22FF}, // Mathematical Operators: U+2200 to U+22FF
		{0x2500, 0x257F}, // Box Drawing: U+2500 to U+257F
	}

	for _, r := range ranges {
		for i := r[0]; i <= r[1]; i++ {
			if mapIdx == threadColorsLength {
				break
			}
			unicodeMap[threadColors[mapIdx].ID] = rune(i)
			mapIdx += 1
		}
	}

	return unicodeMap
}

func ColorAtoi(s string) uint8 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("error converting string to int: %s", err)
		return 0
	}
	return uint8(i)
}

func LoadThreadColors(filePath string) ([]common.ThreadColor, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var threadImg []common.ThreadColor
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

		ThreadColor := common.ThreadColor{
			ID:    id,
			Name:  name,
			Color: color.RGBA{R: ColorAtoi(parts[lineIdx-1]), G: ColorAtoi(parts[lineIdx]), B: ColorAtoi(parts[lineIdx+1])},
		}

		threadImg = append(threadImg, ThreadColor)
	}

	dmcMap := createUnicodeCharMap(threadImg)
	for i := range threadImg {
		threadImg[i].Symbol = string(dmcMap[threadImg[i].ID])
	}

	// dmcMapLength := len(dmcMap)

	// for i := 0; i <= dmcMapLength; i++ {
	// 	threadImg[i].Symbol = string(dmcMap[i])
	// }

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	return threadImg, err
}

func ReduceColors(img image.Image, palette []common.ThreadColor) image.Image {
	bounds := img.Bounds()
	reducedImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 1 {
			originalColor := img.At(x, y)
			nearestColor := colormath.NearestColor(originalColor, palette)
			reducedImg.Set(x, y, nearestColor.Color)
		}
	}

	return reducedImg
}
