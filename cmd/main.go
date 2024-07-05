package main

import (
	"fmt"
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/imageprocessing"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Cross Stitch Image Generator")

	// Image processing UI components
	label := widget.NewLabel("Select a folder to upload an image:")
	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Enter new height")

	numColoursEntry := widget.NewEntry()
	numColoursEntry.SetPlaceHolder("Number of Colours")

	// Image display canvas
	imageCanvas := canvas.NewImageFromImage(nil)
	imageCanvas.FillMode = canvas.ImageFillOriginal

	var currentImage image.Image
	uploadButton := widget.NewButton("Select Folder", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uri == nil {
				return
			}

			// Create a file dialog to select the image file within the chosen folder
			fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				if reader == nil {
					return
				}
				defer reader.Close()

				img, _, err := image.Decode(reader)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}

				// Get the new height from the entry
				newHeight, err := strconv.Atoi(heightEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("invalid height: %w", err), myWindow)
					return
				}

				// Resize the image
				currentImage = imageprocessing.ResizeImage(img, newHeight)

				// Display the resized image on the canvas
				imageCanvas.Image = currentImage
				imageCanvas.Refresh()
			}, myWindow)
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
			fileDialog.SetLocation(uri) // Set the location to the selected folder
			fileDialog.Show()
		}, myWindow).Show()
	})

	legend := widget.NewLabel("Color Legend:")

	// Generate button
	generateButton := widget.NewButton("Generate", func() {
		imgHeight, err := strconv.Atoi(heightEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid height value"), myWindow)
			return
		}
		numColours, err := strconv.Atoi(numColoursEntry.Text)
		if err != nil || numColours < 10 || numColours > 50 {
			dialog.ShowError(fmt.Errorf("Number of colours must be between 10 and 50"), myWindow)
			return
		}

		if currentImage == nil {
			dialog.ShowError(fmt.Errorf("No image loaded"), myWindow)
			return
		}

		resizedImg := imageprocessing.ResizeImage(currentImage, imgHeight)
		threadColours, err := imageprocessing.LoadThreadColours("assets/thread_colours.txt")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to load thread colours"), myWindow)
			return
		}

		threadPalette := imageprocessing.GetPartialPalette(resizedImg, threadColours, numColours)
		reducedImg := imageprocessing.ReduceColors(resizedImg, threadPalette)

		fmt.Println(threadPalette)

		// Ensure the resized and color-reduced image is displayed correctly
		imageCanvas.Image = reducedImg
		imageCanvas.Refresh()

		currentImage = reducedImg

		dialog.ShowInformation("Success", "Image processed successfully", myWindow)
	})

	// Set the content of the window
	myWindow.SetContent(container.NewVBox(
		label,
		heightEntry,
		numColoursEntry,
		uploadButton,
		generateButton,
		imageCanvas,
		legend,
	))

	myWindow.Resize(fyne.NewSize(1400, 950))
	myWindow.ShowAndRun()
}
