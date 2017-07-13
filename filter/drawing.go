package filter

import (
	"image"
	"image/color"
	_ "image/png"
	"math"
)

type Drawing struct {

	data []uint8
	Height int
	Width int
}

func NewDrawing(width, height int, opaque bool) *Drawing {

	var drawing = Drawing{data: make([]uint8, width * height * 4), Height: height, Width: width}

	if opaque {

		for x := 0; x < drawing.Width; x++ {

			for y := 0; y < drawing.Height; y++ {

				drawing.Set(x, y, 0, 0, 0, math.MaxUint8)
			}
		}

	}

	return &drawing
}

func DrawingFromImage(source image.Image) *Drawing {

	var width = source.Bounds().Max.X
	var height = source.Bounds().Max.Y

	var result = Drawing{data: make([]uint8, width * height * 4), Height: height, Width: width}

	for x := 0; x < source.Bounds().Max.X; x++ {

		for y := 0; y < source.Bounds().Max.Y; y++ {

			result.SetColor(x, y, source.At(x, y))
		}
	}

	return &result
}

func DrawingToImage(source *Drawing) image.Image {

	var result = image.NewRGBA64(image.Rect(0, 0, source.Width, source.Height))

	for x := 0; x < source.Width; x++ {

		for y := 0; y < source.Height; y++ {

			result.Set(x, y, source.GetColor(x, y))
		}
	}

	return result
}

func (img *Drawing) SetColor(x, y int, c color.Color) {

	var r, g, b, a = c.RGBA()

	img.Set(x, y, uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8))
}

func (img *Drawing) Set(x, y int, r, g, b, a uint8) {

	var i = img.indexOf(x, y)
	img.data[i] = r
	img.data[i + 1] = g
	img.data[i + 2] = b
	img.data[i + 3] = a
}

func (img *Drawing) GetColor(x, y int) color.Color {

	var r, g, b, a = img.Get(x, y)
	r64 := uint16(r) << 8 | uint16(r)
	g64 := uint16(g) << 8 | uint16(g)
	b64 := uint16(b) << 8 | uint16(b)
	a64 := uint16(a) << 8 | uint16(a)

	return &color.RGBA64{R: r64, G: g64, B: b64, A: a64}
}

func (img *Drawing) Get(x, y int) (r, g, b, a uint8){

	var i = img.indexOf(x, y)
	r = img.data[i]
	g = img.data[i + 1]
	b = img.data[i + 2]
	a = img.data[i + 3]
	return
}

func (img *Drawing) indexOf(x int, y int) int {

	return (y * img.Width + x) * 4
}

func (img *Drawing) Opaque() bool {

	for x := 0; x < img.Width; x++ {

		for y := 0; y < img.Height; y++ {

			var i = img.indexOf(x, y)
			if img.data[i + 3] != math.MaxUint8 {

				return false
			}
		}
	}

	return true
}