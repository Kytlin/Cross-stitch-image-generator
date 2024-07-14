package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/imageprocessing"
)

var threadPalette []common.ThreadColour
var rectangles [][]*canvas.Rectangle

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
	//unicodeLabel := widget.NewLabel("Unicode Symbols: \u2764\u2600\u2601") // Example symbols

	// Apply custom font
	customFontResource := fyne.NewStaticResource("CustomFont", customFont)
	fyne.CurrentApp().Settings().SetTheme(&myTheme{font: customFontResource})

	// Image processing UI components
	label := widget.NewLabel("Select a folder to upload an image:")

	// HEIGHT
	defaultHeight := 50.0
	heightValue := binding.NewFloat()
	heightValue.Set(defaultHeight)
	heightSlider := widget.NewSliderWithData(20.0, 200.0, heightValue)

	// Custom label to display integer value
	heightLabel := widget.NewLabelWithData(binding.NewString())
	heightValue.AddListener(binding.NewDataListener(func() {
		floatVal, _ := heightValue.Get()
		intVal := int(floatVal)
		heightLabel.SetText(strconv.Itoa(intVal))
	}))

	// NUM COLOURS
	defaultNumColours := 50.0
	numColours := binding.NewFloat()
	numColours.Set(defaultNumColours)
	numColoursSlider := widget.NewSliderWithData(10.0, 200.0, numColours)

	// Custom label to display integer value
	numColoursLabel := widget.NewLabelWithData(binding.NewString())
	numColours.AddListener(binding.NewDataListener(func() {
		floatVal, _ := numColours.Get()
		intVal := int(floatVal)
		numColoursLabel.SetText(strconv.Itoa(intVal))
	}))

	// Image display canvas
	imageCanvas := canvas.NewImageFromImage(nil)
	imageCanvas.FillMode = canvas.ImageFillOriginal

	legend := getLegend()
	uploadButton, generateButton := getUploadAndGenerateButtons(heightLabel, numColoursLabel, legend, myWindow, imageCanvas)

	myWindow.SetContent(container.NewVBox(
		label,
		heightLabel,
		heightSlider,
		numColoursLabel,
		numColoursSlider,
		uploadButton,
		generateButton,
		imageCanvas,
		legend,
		// unicodeLabel,
	))

	myWindow.Resize(fyne.NewSize(1400, 800))
	myWindow.ShowAndRun()
}

func getUploadAndGenerateButtons(heightEntry *widget.Label, numColoursEntry *widget.Label, legend fyne.CanvasObject, myWindow fyne.Window, imageCanvas *canvas.Image) (fyne.CanvasObject, fyne.CanvasObject) {
	var currentImage image.Image

	currentDir, _ := os.Getwd()
	curUri := storage.NewFileURI(currentDir)
	uri, _ := storage.ListerForURI(curUri)
	_ = curUri

	// Upload
	uploadButton := widget.NewButton("Select Folder", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			fmt.Println("Selected file:", reader.URI().Path())
			imagePath := reader.URI().Path()
			imageFile, err := os.Open(imagePath)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			img, _, err := image.Decode(imageFile)
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
	})

	// Generate
	generateButton := widget.NewButton("Generate", func() {
		imgHeight, err := strconv.Atoi(heightEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid height value"), myWindow)
			return
		}
		numColours, err := strconv.Atoi(numColoursEntry.Text)
		if err != nil || numColours < 10 || numColours > 200 {
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

		threadPalette = imageprocessing.GetPartialPalette(resizedImg, threadColours, numColours)
		reducedImg := imageprocessing.ReduceColors(resizedImg, threadPalette)

		fmt.Println(threadPalette)

		// Ensure the resized and color-reduced image is displayed correctly
		imageCanvas.Image = resizedImg

		colourGrid := imageprocessing.GenerateColourGrid(reducedImg, []common.ThreadColour{}, false)
		gridImage := generateImageFromGrid(colourGrid, false, false)
		updateGrid(colourGrid)

		imageCanvas.Image = gridImage

		imageCanvas.Refresh()
		legend = getLegend()
		legend.Refresh()

		currentImage = reducedImg

		dialog.ShowInformation("Success", "Image processed successfully", myWindow)
	})

	return uploadButton, generateButton
}

func getLegend() fyne.CanvasObject {
	legend := widget.NewTable(
		func() (int, int) {
			// Returning the number of rows and columns
			if threadPalette == nil {
				return 0, 4
			}
			return len(threadPalette), 4
		},
		func() fyne.CanvasObject {
			// Create a new label for each cell
			return container.NewVBox(widget.NewLabel(""), canvas.NewRectangle(color.Black))
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*fyne.Container).Objects[0].(*widget.Label)
			i := o.(*fyne.Container).Objects[1].(*canvas.Rectangle)
			l.Show()
			i.Hide()

			if threadPalette != nil && id.Row < len(threadPalette) {
				switch id.Col {
				case 0:
					l.SetText("[symbol]")
				case 1:
					l.SetText("DMC " + strconv.Itoa(threadPalette[id.Row].ID))
				case 2:
					l.SetText(threadPalette[id.Row].Name)
				case 3:
					l.Hide()
					colour := threadPalette[id.Row].Colour
					colour.A = 255
					i.FillColor = colour
					i.SetMinSize(fyne.NewSize(20, 20))
					i.Show()
				}
			}
		})

	// Set column widths
	legend.SetColumnWidth(0, 80)
	legend.SetColumnWidth(1, 120)
	legend.SetColumnWidth(2, 450)

	// Create header for the table
	legend.CreateHeader = func() fyne.CanvasObject {
		return widget.NewButton("", func() {})
	}
	legend.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
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

	// Scroll container to make table height adjustable
	scrollContainer := container.NewScroll(legend)
	scrollContainer.SetMinSize(fyne.NewSize(650, 300))

	return scrollContainer
}

func generateImageFromGrid(grid [][]common.ThreadColour, showSymbol bool, useStich bool) image.Image {
	numRows := len(grid)
	numCols := len(grid[0])

	cellSize := 20
	imgWidth := numCols * cellSize
	imgHeight := numRows * cellSize

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			cell := grid[row][col]
			x := col * cellSize
			y := row * cellSize

			cellColor := cell.Colour
			cellColor.A = 255

			if useStich {
				stitchThickness := 3

				for i := 0; i < cellSize; i++ {
					for t := 0; t < stitchThickness; t++ {
						img.Set(x+i, y+i+t, cellColor)            // top left diagonal
						img.Set(x+i, y+cellSize-1-i-t, cellColor) // bottom left diagonal
						img.Set(x+i+t, y+i, cellColor)            // top left diagonal (offset)
						img.Set(x+i+t, y+cellSize-1-i, cellColor) // bottom left diagonal (offset)
					}
				}
			} else {
				// Fill with color
				draw.Draw(img, image.Rect(x, y, x+cellSize, y+cellSize), &image.Uniform{cellColor}, image.Point{}, draw.Src)
			}

			if showSymbol {
				// fnt, _ := opentype.Parse(customFont)
				// face, _ := opentype.NewFace(fnt, &opentype.FaceOptions{
				//  Size:    float64(cellSize),
				//  DPI:     72,
				//  Hinting: font.HintingFull,
				// })
				// defer face.Close()
				// drawer := &font.Drawer{
				//  Dst: img,
				//  Src: image.White,
				//  Face: face,
				// }
				// drawer.Dot = fixed.Point26_6{
				//  X: fixed.I(x + cellSize/4),
				//  Y: fixed.I(y + cellSize - cellSize/4),
				// }
				// symbol := string(cell.Symbol)
				// drawer.DrawString(symbol)
			}

			// Border around cell
			borderColor := color.Black
			borderThickness := 1
			draw.Draw(img, image.Rect(x, y, x+cellSize, y+borderThickness), &image.Uniform{borderColor}, image.Point{}, draw.Src)
			draw.Draw(img, image.Rect(x, y, x+borderThickness, y+cellSize), &image.Uniform{borderColor}, image.Point{}, draw.Src)
			draw.Draw(img, image.Rect(x, y+cellSize-borderThickness, x+cellSize, y+cellSize), &image.Uniform{borderColor}, image.Point{}, draw.Src)
			draw.Draw(img, image.Rect(x+cellSize-borderThickness, y, x+cellSize, y+cellSize), &image.Uniform{borderColor}, image.Point{}, draw.Src)
		}
	}

	return img
}

func updateRectangleColor(row, col int, threadColour common.ThreadColour) {
	if row >= 0 && row < len(rectangles) && col >= 0 && col < len(rectangles[row]) {
		r, g, b, _ := threadColour.Colour.RGBA()
		rectangles[row][col].FillColor = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		rectangles[row][col].Refresh() // Refresh to apply the color change
	}
}

func updateGrid(colours [][]common.ThreadColour) {
	for i := 0; i < len(colours); i++ {
		for j := 0; j < len(colours[i]); j++ {
			updateRectangleColor(i, j, colours[i][j])
		}
	}
}
