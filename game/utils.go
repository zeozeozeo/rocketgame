package game

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadRGBAImage(f io.Reader) (*image.RGBA, error) {
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	// convert to image.RGBA
	b := src.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	return m, nil
}

func LoadEbitenImage(data []byte) *ebiten.Image {
	img, err := LoadRGBAImage(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
