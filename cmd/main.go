package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
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
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/Kytlin/Cross-stitch-image-generator/pkg/common"
	"github.com/Kytlin/Cross-stitch-image-generator/pkg/imageprocessing"
)

var currentImage image.Image

var showSymbol, useStitch bool

var threadPalette []common.ThreadColor
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

	// Apply custom font
	customFontResource := fyne.NewStaticResource("CustomFont", customFont)
	fyne.CurrentApp().Settings().SetTheme(&myTheme{font: customFontResource})

	// Image processing UI components
	label := widget.NewLabel("Select a folder to upload an image:")

	// HEIGHT
	defaultHeight := 30.0
	heightValue := binding.NewFloat()
	heightValue.Set(defaultHeight)
	heightSlider := widget.NewSliderWithData(10.0, 50.0, heightValue)

	// Custom label to display integer value
	heightLabel := widget.NewLabelWithData(binding.NewString())
	heightValue.AddListener(binding.NewDataListener(func() {
		floatVal, _ := heightValue.Get()
		intVal := int(floatVal)
		heightLabel.SetText("Height:\t " + strconv.Itoa(intVal))
	}))

	// NUM colorS
	defaultNumColors := 30.0
	numColors := binding.NewFloat()
	numColors.Set(defaultNumColors)
	numColorsSlider := widget.NewSliderWithData(10.0, 50.0, numColors)

	// Custom label to display integer value
	numColorsLabel := widget.NewLabelWithData(binding.NewString())
	numColors.AddListener(binding.NewDataListener(func() {
		floatVal, _ := numColors.Get()
		intVal := int(floatVal)
		numColorsLabel.SetText("Number of Thread Colors: " + strconv.Itoa(intVal))
	}))

	gridDownloadChoice := widget.NewRadioGroup([]string{"Filled color and symbol", "Filled color", "X stitch"}, func(value string) {
		if value == "Filled color and symbol" {
			showSymbol = true
			useStitch = false
		} else if value == "X stitch" {
			showSymbol = false
			useStitch = true
		} else {
			showSymbol = false
			useStitch = false
		}
	})

	// Image display canvas
	imageCanvas := canvas.NewImageFromImage(nil)
	imageCanvas.FillMode = canvas.ImageFillOriginal

	// Create a placeholder for the legend table
	legend := getLegend()
	uploadButton, resizeButton, generateButton := getUploadAndGenerateButtons(heightSlider, numColorsSlider, myWindow, imageCanvas, customFont)

	myWindow.SetContent(container.NewVBox(
		label,
		heightLabel,
		heightSlider,
		numColorsLabel,
		numColorsSlider,
		gridDownloadChoice,
		uploadButton,
		resizeButton,
		generateButton,
		imageCanvas,
		legend,
	))

	myWindow.Resize(fyne.NewSize(1400, 800))
	myWindow.ShowAndRun()
}

func getUploadAndGenerateButtons(heightSlider *widget.Slider, numColorsSlider *widget.Slider, myWindow fyne.Window, imageCanvas *canvas.Image, customFont []byte) (fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject) {
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

			decodedImg, _, err := image.Decode(imageFile)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			currentImage = decodedImg

			// Process image based on height input
			processedImage := imageprocessing.ResizeImage(currentImage, int(heightSlider.Value))

			// Display image on canvas
			colorGrid := imageprocessing.GenerateColorGrid(processedImage, []common.ThreadColor{}, false)
			gridImage := generateImageFromGrid(colorGrid, false, false, customFont)

			imageCanvas.Image = gridImage
			imageCanvas.Refresh()
		}, myWindow)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fileDialog.SetLocation(uri)
		fileDialog.Show()
	})

	// Resize
	resizeButton := widget.NewButton("Resize Image", func() {
		if currentImage == nil {
			dialog.ShowError(fmt.Errorf("No image loaded"), myWindow)
			return
		}

		// Resize the image
		resizedImage := imageprocessing.ResizeImage(currentImage, int(heightSlider.Value))

		// Display image on canvas
		colorGrid := imageprocessing.GenerateColorGrid(resizedImage, []common.ThreadColor{}, false)
		gridImage := generateImageFromGrid(colorGrid, showSymbol, useStitch, customFont)

		imageCanvas.Image = gridImage
		imageCanvas.Refresh()
	})

	// Generate
	generateButton := widget.NewButton("Generate", func() {
		imgHeight := heightSlider.Value
		numColors := numColorsSlider.Value

		if currentImage == nil {
			dialog.ShowError(fmt.Errorf("No image loaded"), myWindow)
			return
		}

		resizedImg := imageprocessing.ResizeImage(currentImage, int(imgHeight))
		threadColors, err := imageprocessing.LoadThreadColors("assets/thread_colors.txt")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to load thread colors"), myWindow)
			return
		}

		threadPalette = imageprocessing.GetPartialPalette(resizedImg, threadColors, int(numColors))
		reducedImg := imageprocessing.ReduceColors(resizedImg, threadPalette)

		// Display the resized and color-reduced image on canvas
		colorGrid := imageprocessing.GenerateColorGrid(reducedImg, threadColors, true)
		gridImage := generateImageFromGrid(colorGrid, showSymbol, useStitch, customFont)
		updateGrid(colorGrid)

		imageCanvas.Image = gridImage

		imageCanvas.Refresh()

		legend := getLegend()
		legend.Refresh()

		// Save the generated images
		saveGeneratedImages(resizedImg, threadColors, customFont, myWindow)

		dialog.ShowInformation("Success", "Image processed and saved successfully", myWindow)
	})

	return uploadButton, resizeButton, generateButton
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
			return container.NewHBox(widget.NewLabel(""), canvas.NewRectangle(color.Black))
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*fyne.Container).Objects[0].(*widget.Label)
			i := o.(*fyne.Container).Objects[1].(*canvas.Rectangle)
			l.Show()
			i.Hide()

			if threadPalette != nil && id.Row < len(threadPalette) {
				switch id.Col {
				case 0:
					l.SetText(threadPalette[id.Row].Symbol)
				case 1:
					l.SetText("DMC " + strconv.Itoa(threadPalette[id.Row].ID))
				case 2:
					l.SetText(threadPalette[id.Row].Name)
				case 3:
					l.Hide()
					color := threadPalette[id.Row].Color
					color.A = 255
					i.FillColor = color
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
		return container.NewHBox(
			widget.NewLabel("Symbol"),
			widget.NewLabel("Number"),
			widget.NewLabel("Name"),
			widget.NewLabel("Color"),
		)
	}

	legend.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		hbox := o.(*fyne.Container)
		label := hbox.Objects[id.Col].(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("Symbol")
		case 1:
			label.SetText("Number")
		case 2:
			label.SetText("Name")
		case 3:
			label.SetText("Color")
		}
	}

	// Scroll container to make table height adjustable
	scrollContainer := container.NewScroll(legend)
	scrollContainer.SetMinSize(fyne.NewSize(650, 300))

	return scrollContainer
}

func generateImageFromGrid(grid [][]common.ThreadColor, showSymbol bool,
	useStitch bool, customFont []byte) image.Image {
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

			cellColor := cell.Color
			cellColor.A = 255

			if useStitch {
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
				fnt, _ := opentype.Parse(customFont)
				face, _ := opentype.NewFace(fnt, &opentype.FaceOptions{
					Size:    float64(cellSize),
					DPI:     72,
					Hinting: font.HintingFull,
				})
				defer face.Close()

				fontColor := image.White
				if (float32(cell.Color.R)*0.299 + float32(cell.Color.G)*0.587 + float32(cell.Color.B)*0.114) > 186 {
					fontColor = image.Black
				}

				drawer := &font.Drawer{
					Dst:  img,
					Src:  fontColor,
					Face: face,
				}
				drawer.Dot = fixed.Point26_6{
					X: fixed.I(x + cellSize/4),
					Y: fixed.I(y + cellSize - cellSize/4),
				}
				symbol := cell.Symbol
				drawer.DrawString(symbol)
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

func saveGeneratedImages(img image.Image, threadColors []common.ThreadColor, customFont []byte, myWindow fyne.Window) {
	basePath := "output/"
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to create output directory"), myWindow)
		return
	}

	types := []struct {
		showSymbol bool
		useStitch  bool
		name       string
	}{
		{true, false, "filled_color_and_symbol.jpg"},
		{false, false, "filled_color.jpg"},
		{false, true, "x_stitch.jpg"},
	}

	for _, t := range types {
		colorGrid := imageprocessing.GenerateColorGrid(img, threadColors, true)
		gridImage := generateImageFromGrid(colorGrid, t.showSymbol, t.useStitch, customFont)

		saveImageToFile(gridImage, basePath+t.name)
	}
}

func saveImageToFile(img image.Image, pathname string) {
	file, err := os.Create(pathname)
	if err != nil {
		fmt.Println("Failed to save image:", err)
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("Failed to encode image:", err)
	}
}

func updateRectangleColor(row, col int, threadColor common.ThreadColor) {
	if row >= 0 && row < len(rectangles) && col >= 0 && col < len(rectangles[row]) {
		r, g, b, _ := threadColor.Color.RGBA()
		rectangles[row][col].FillColor = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		rectangles[row][col].Refresh() // Refresh to apply the color change
	}
}

func updateGrid(colors [][]common.ThreadColor) {
	for i := 0; i < len(colors); i++ {
		for j := 0; j < len(colors[i]); j++ {
			updateRectangleColor(i, j, colors[i][j])
		}
	}
}
