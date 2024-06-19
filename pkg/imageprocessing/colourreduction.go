package imageprocessing

import (
	"bufio"
	"image/color"
	"log"
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
