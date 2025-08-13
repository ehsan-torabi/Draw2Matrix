// Package main provides image processing functionality for the Draw2Matrix application
package main

import (
	"fyne.io/fyne/v2"
	"github.com/anthonynsimon/bild/transform"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

// imageProcessor converts a raw image into a processed grayscale image
// It crops, scales, and binarizes the image according to the specified dimensions
func imageProcessor(img image.Image, size fyne.Size, position fyne.Position) *image.Gray {
	// Crop the image with padding
	//croppedImg, err := cutter.Crop(img, cutter.Config{
	//	Width:  int(size.Width),  // Add padding to width
	//	Height: int(size.Height), // Add padding to height
	//	Anchor: image.Point{X: int(position.X), Y: int(position.Y)},
	//})
	//if err != nil {
	//	panic(err)
	//}

	// Create a black and white palette image
	dst := image.NewPaletted(
		image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y),
		color.Palette{color.White, color.Black},
	)
	draw.Draw(dst, img.Bounds(), img, image.Point{}, draw.Over)

	// Create final grayscale image with desired dimensions
	final := image.NewGray(image.Rect(0, 0, Options.MatrixRow, Options.MatrixCol))
	draw.CatmullRom.Scale(final, final.Rect, img, img.Bounds(), draw.Over, nil)

	// Binarize the image (convert to pure black and white)
	for i := 0; i != Options.MatrixRow; i++ {
		for j := 0; j != Options.MatrixCol-1; j++ {
			if final.GrayAt(i, j).Y < 255 {
				final.Set(i, j, color.Black)
			} else {
				final.Set(i, j, color.White)
			}
		}
	}

	return final
}

// image2BinaryMatrix converts a grayscale image to a binary matrix
// The image is rotated and flipped to match the desired orientation
// Black pixels (0) are converted to 1, white pixels (255) are converted to 0
func image2BinaryMatrix(img *image.Gray) [][]int8 {
	// Rotate and flip the image for correct orientation
	temp := transform.Rotate(img, 90, nil)
	temp = transform.FlipH(temp)

	// Draw the transformed image back to the original
	rect := image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())
	draw.Draw(img, rect, temp, image.Point{}, draw.Over)

	// Create binary matrix
	result := make([][]int8, Options.MatrixRow-1)
	for i := 0; i != Options.MatrixRow-1; i++ {
		result[i] = make([]int8, Options.MatrixCol-1)
		for j := 0; j != Options.MatrixCol-1; j++ {
			// Convert black to 1, white to 0
			if img.GrayAt(i, j).Y == 0 {
				result[i][j] = int8(1)
			} else {
				result[i][j] = int8(0)
			}
		}
	}
	return result
}

// captureAndProcessImage captures the current window content and processes it
// Returns a grayscale image that has been cropped and processed
func captureAndProcessImage(w fyne.Window, p *PaintWidget) *image.Gray {
	img := w.Canvas().Capture()
	return imageProcessor(img, p.Size(), p.Position())
}
