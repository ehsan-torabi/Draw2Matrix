package main

import (
	"fyne.io/fyne/v2"
	"github.com/anthonynsimon/bild/transform"
	"github.com/oliamb/cutter"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

const (
	MATRIX_COL_NUM = 20
	MATRIX_ROW_NUM
)

func imageProcessor(img image.Image, size fyne.Size, position fyne.Position) *image.Gray {
	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  int(size.Width) + 100,
		Height: int(size.Height) + 65,
		Anchor: image.Point{X: int(position.X) + 10, Y: int(position.Y) + 10},
	})
	if err != nil {
		panic(err)
	}
	dst := image.NewPaletted(image.Rect(0, 0, croppedImg.Bounds().Size().X, croppedImg.Bounds().Size().Y), color.Palette{color.White, color.Black})
	draw.Draw(dst, img.Bounds(), croppedImg, image.Point{}, draw.Over)
	final := image.NewGray(image.Rect(0, 0, MATRIX_ROW_NUM, MATRIX_COL_NUM))
	draw.CatmullRom.Scale(final, final.Rect, croppedImg, croppedImg.Bounds(), draw.Over, nil)
	for i := 1; i != MATRIX_ROW_NUM; i++ {
		for j := 1; j != MATRIX_COL_NUM-1; j++ {
			y := final.GrayAt(i, j).Y
			if 255 > y {
				final.Set(i, j, color.Black)
			} else {
				final.Set(i, j, color.White)
			}

		}
	}
	return final
}

func toBinaryMatrix(img *image.Gray) [][]int8 {
	temp := transform.Rotate(img, 90, nil)
	temp = transform.FlipH(temp)
	rect := image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())
	draw.Draw(img, rect, temp, image.Point{}, draw.Over)
	var result [][]int8
	result = make([][]int8, MATRIX_ROW_NUM-1)
	for i := 0; i != MATRIX_ROW_NUM-1; i++ {
		result[i] = make([]int8, MATRIX_COL_NUM-1)
		for j := 0; j != MATRIX_COL_NUM-1; j++ {
			y := img.GrayAt(i, j).Y
			if y == 0 {
				result[i][j] = int8(1)
			} else {
				result[i][j] = int8(0)
			}

		}
	}
	return result
}

func captureAndProcessImage(w fyne.Window, p *PaintWidget) *image.Gray {
	img := w.Canvas().Capture()
	result := imageProcessor(img, p.Size(), p.Position())
	return result
}
