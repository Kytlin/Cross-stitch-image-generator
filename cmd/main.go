package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/imageprocessing"
)

var threadPalette []common.ThreadColour

// loadCustomFont reads and loads a TTF font from the given path.
func loadCustomFont(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fontBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fontBytes, nil
}

// myTheme is a custom Fyne theme that uses a custom font.
type myTheme struct {
	font fyne.Resource
}

func (m *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return m.font
}

func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Cross Stitch Image Generator")

	fontPath := "assets/DejaVuSans.ttf"
	customFont, err := loadCustomFont(fontPath)
	if err != nil {
		fmt.Println("Failed to load custom font: %v", err)
	}

	// Create a label with Unicode symbols
	unicodeLabel := widget.NewLabel("Unicode Symbols: \u2764\u2600\u2601") // Example symbols

	// Apply custom font
	customFontResource := fyne.NewStaticResource("CustomFont", customFont)
	fyne.CurrentApp().Settings().SetTheme(&myTheme{font: customFontResource})

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

	legend := widget.NewTableWithHeaders(func() (int, int) {
		return len(threadPalette), 3
	},
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			return l
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Truncation = fyne.TextTruncateEllipsis
			switch id.Col {
			case 0:
				// This uses a lookup table
				l.SetText("[symbol]")
			case 1:
				l.SetText("DMC " + strconv.Itoa(threadPalette[id.Row].ID))
			case 2:
				l.SetText(threadPalette[id.Row].Name)
			}
		})
	legend.SetColumnWidth(0, 80)
	legend.SetColumnWidth(1, 120)
	legend.SetColumnWidth(2, 450)

	legend.CreateHeader = func() fyne.CanvasObject {
		return widget.NewButton("000", func() {})
	}
	legend.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("Symbol")
			case 1:
				b.SetText("Number")
			case 2:
				b.SetText("Name")
			}
			b.Refresh()
		}
	}

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

		legend.Refresh()

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
		unicodeLabel,
	))

	myWindow.Resize(fyne.NewSize(1400, 950))
	myWindow.ShowAndRun()
}
