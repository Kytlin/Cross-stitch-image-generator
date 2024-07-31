# Cross Stitch Generator

This desktop app was created using Go Fyne library. Run `go run cmd/main.go` to start app. Here is an example image followed by its cross stitched versions generated from app:

<br>

<center>
    <img src="lavender_field.jpeg" style="width:50%">
</center>

<br>

<p float="left">
    <img src="example_images/filled_color.jpg" style="width:30%;padding:0.5em">
    <img src="example_images/filled_color_and_symbol.jpg" style="width:30%;padding:0.5em">
    <img src="example_images/x_stitch.jpg" style="width:30%;padding:0.5em">
</p>

First, the app needs to know which image you want. This can be done by clicking "Select Folder" button and navigating through internal file explorer of one's local machine (image file types supported are `jpeg`/`jpg` or `png`). To resize image, one adjusts the "Height" slider and clicks "Resize" button; width is automatically determined based on image aspect ratio. One can adjust "Number of Color" slider to display the image as a grid using the specified amount of colors, which is accomplished by [color quantization](https://en.wikipedia.org/wiki/Color_quantization#:~:text=In%20computer%20graphics%2C%20color%20quantization,possible%20to%20the%20original%20image.). 

After clicking "Generate" button with an option selected in radio menu, a cross stitch board will display with the corresponding legend. The options are 

- *Filled color -* default option 
- *Filled color with symbols -* add symbols onto image to aid stitching the right color
- *X-stitch -* add visual appeal to default by replacing each cell as a cross stitch pattern

## Demos

The following will demonstrate the GUI accessibility via GIF below:

<br>

An output folder is created with images for each option shown in GUI

<center>
    <img src="generate_image_and_show_output.gif" width="90%" style="margin:2rem; border-radius:1rem">
</center>

A user can always changes images within the program before generating.

<center>
    <img src="select_diff_image_and_generate.gif" width="90%" style="margin:2rem; border-radius:1rem">
</center>

One can change inputs in sliders before generating.

<center>
    <img src="resize_image.gif" width="90%" style="margin:2rem; border-radius:1rem">
</center>

Any of the three options can be applied after each click on "Generate" button.

<center>
    <img src="generate_per_image_option.gif" width="90%" style="margin:2rem; border-radius:1rem">
</center>

Note that the legend updates according to each cross stitch generation.