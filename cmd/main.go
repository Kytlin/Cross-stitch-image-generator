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

	colorCountEntry := widget.NewEntry()
	colorCountEntry.SetPlaceHolder("Enter number of colors")

	// Image display canvas
	imageCanvas := canvas.NewImageFromImage(nil)
	imageCanvas.FillMode = canvas.ImageFillOriginal

	var resizedImage image.Image
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
				resizedImg := imageprocessing.ResizeImage(img, newHeight)

				// Display the resized image on the canvas
				imageCanvas.Image = resizedImg
				imageCanvas.Refresh()

				// Store the resized image for further processing
				resizedImage = resizedImg
			}, myWindow)
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
			fileDialog.SetLocation(uri) // Set the location to the selected folder
			fileDialog.Show()
		}, myWindow).Show()
	})

	legend := widget.NewLabel("Color Legend:")

	generateButton := widget.NewButton("Generate", func() {
		if resizedImage == nil {
			dialog.ShowError(fmt.Errorf("no image selected"), myWindow)
			return
		}

		// Get the number of colors from the entry
		numColors, err := strconv.Atoi(colorCountEntry.Text)
		if err != nil || numColors <= 0 {
			dialog.ShowError(fmt.Errorf("invalid number of colors: %w", err), myWindow)
			return
		}

		// threadColors := imageprocessing.getPartialPalette(resizedImage, numColors) [TODO]
	})

	// Set the content of the window
	myWindow.SetContent(container.NewVBox(
		label,
		heightEntry,
		colorCountEntry,
		uploadButton,
		generateButton,
		imageCanvas,
		legend,
	))

	myWindow.Resize(fyne.NewSize(1400, 950))
	myWindow.ShowAndRun()
}
