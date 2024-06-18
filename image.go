package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/image/draw"
)

// LoadImage loads an image from the specified file path.
func LoadImage(filePath string) (image.Image, error) {
	// Open the file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file extension to determine the decoder to use.
	ext := filepath.Ext(filePath)
	var img image.Image

	// Decode the image based on the file extension.
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file) // outputs a string of rgb
		if err != nil {
			return nil, err
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	return img, nil
}

// SaveImage saves an image to the specified file path with the given format.
func SaveImage(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, nil)
	case ".png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

func ResizeImage(img image.Image, newHeight int) image.Image {
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dx()

	// Calculate the new width to maintain the aspect ratio
	newWidth := (newHeight * imgWidth) / imgHeight

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func main() {
	// Ensure correct number of arguments
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: go run main.go resize [image height] [image] [input file name] [output file name (optional)]")
		os.Exit(1)
	}

	// Parse command-line arguments
	imgOption := os.Args[1]
	imgHeight, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: Invalid height. Please provide a valid integer.")
		os.Exit(1)
	}
	inputFilePath := os.Args[3]

	// Default output file path
	var outputFilePath string
	if len(os.Args) == 5 {
		outputFilePath = os.Args[4]
	} else {
		// Generate canonical output file name based on input file name
		baseName := filepath.Base(inputFilePath)
		ext := filepath.Ext(baseName)
		outputFilePath = "resized_" + baseName[:len(baseName)-len(ext)] + ext
	}

	// Validate operation
	if imgOption != "resize" {
		fmt.Println("Error: Unsupported operation. Only 'resize' is supported.")
		os.Exit(1)
	}

	// Load the image
	img, err := LoadImage(inputFilePath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		os.Exit(1)
	}

	// Resize the image
	resizedImg := ResizeImage(img, imgHeight)

	// Save the resized image
	err = SaveImage(outputFilePath, resizedImg)
	if err != nil {
		fmt.Println("Error saving image:", err)
		os.Exit(1)
	}

	fmt.Println("Image resized and saved successfully to", outputFilePath)
}
