package imageprocessing

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
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

// minThreadColorFields is the smallest usable line: an id, a single name field,
// then the trailing R, G, B and hex values.
const minThreadColorFields = 6

func parseColorComponent(s string) (uint8, error) {
	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid color component %q: %w", s, err)
	}
	return uint8(i), nil
}

// LoadThreadColors reads a tab separated thread list, one thread per line, in
// the form "id<TAB>name<TAB>R<TAB>G<TAB>B<TAB>hex". A name may itself span
// several tab separated fields, so the colors are read from the end of the line.
func LoadThreadColors(filePath string) ([]common.ThreadColor, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var threadImg []common.ThreadColor
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '\t'
		})
		if len(parts) < minThreadColorFields {
			return nil, fmt.Errorf("%s line %d: expected at least %d tab separated fields, got %d",
				filePath, lineNum, minThreadColorFields, len(parts))
		}

		// The last four fields are R, G, B and hex, so the name is everything
		// between the id and those.
		redIdx := len(parts) - 4

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("%s line %d: invalid thread id %q: %w", filePath, lineNum, parts[0], err)
		}

		var rgb [3]uint8
		for i := range rgb {
			rgb[i], err = parseColorComponent(parts[redIdx+i])
			if err != nil {
				return nil, fmt.Errorf("%s line %d: %w", filePath, lineNum, err)
			}
		}

		threadImg = append(threadImg, common.ThreadColor{
			ID:    id,
			Name:  strings.Join(parts[1:redIdx], " "),
			Color: color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2]},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", filePath, err)
	}

	dmcMap := createUnicodeCharMap(threadImg)
	for i := range threadImg {
		threadImg[i].Symbol = string(dmcMap[threadImg[i].ID])
	}

	return threadImg, nil
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
