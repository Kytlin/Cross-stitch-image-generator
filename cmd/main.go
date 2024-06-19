package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/imageprocessing"
)

// validateArguments ensures the correct number of arguments.
func validateArguments() (string, int, string, string) {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: go run main.go resize [height] [input_image] [output_image (optional)]")
		os.Exit(1)
	}

	imgOption := os.Args[1]
	imgHeight, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: Invalid height. Please provide a valid integer.")
		os.Exit(1)
	}
	inputFilePath := os.Args[3]

	var outputFilePath string
	if len(os.Args) == 5 {
		outputFilePath = os.Args[4]
	} else {
		baseName := filepath.Base(inputFilePath)
		ext := filepath.Ext(baseName)
		outputFilePath = "resized_" + baseName[:len(baseName)-len(ext)] + ext
	}

	return imgOption, imgHeight, inputFilePath, outputFilePath
}

func main() {
	imgOption, imgHeight, inputFilePath, outputFilePath := validateArguments()

	if imgOption != "resize" {
		fmt.Println("Error: Unsupported operation. Only 'resize' is supported.")
		os.Exit(1)
	}

	img, err := imageprocessing.LoadImage(inputFilePath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		os.Exit(1)
	}

	resizedImg := imageprocessing.ResizeImage(img, imgHeight)

	err = imageprocessing.SaveImage(outputFilePath, resizedImg)
	if err != nil {
		fmt.Println("Error saving image:", err)
		os.Exit(1)
	}

	fmt.Println("Image resized and saved successfully to", outputFilePath)

	threadColours, err := imageprocessing.LoadThreadColours("assets/thread_colours.txt")
	if err != nil {
		fmt.Println("Error loading thread colours:", err)
		os.Exit(1)
	}
	fmt.Println("Thread colours loaded successfully")

	reducedImg := imageprocessing.ReduceColors(resizedImg, threadColours)

	reducedOutputFilePath := "reduced_" + filepath.Base(outputFilePath)
	err = imageprocessing.SaveImage(reducedOutputFilePath, reducedImg)
	if err != nil {
		fmt.Println("Error saving reduced color image:", err)
		os.Exit(1)
	}

	fmt.Println("Image resized, color-reduced, and saved successfully to", reducedOutputFilePath)

	imageprocessing.GenerateGrid(reducedImg, threadColours)
}
