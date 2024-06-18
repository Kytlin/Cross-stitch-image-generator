package imageprocessing

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

// LoadImage loads an image from the specified file path.
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
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

// ResizeImage resizes the image to the specified height while maintaining aspect ratio.
func ResizeImage(img image.Image, newHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth := (newHeight * width) / height

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
