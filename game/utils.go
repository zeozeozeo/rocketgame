package game

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"math"

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

func Lerp(a, b, t float64) float64 {
	return a*(1.0-t) + (b * t)
}

func fposmod(a, b float64) float64 {
	if a >= 0 {
		return math.Mod(a, b)
	}
	return b - math.Mod(-a, b)
}

func normalizeAngle(x float64) float64 {
	return fposmod(x+math.Pi, 2.0*math.Pi) - math.Pi
}

func LerpAngle(a, b, t float64) float64 {
	
	if math.Abs(a-b) >= math.Pi {
		if a > b {
			a = normalizeAngle(a) - 2.0*math.Pi
		} else {
			b = normalizeAngle(b) - 2.0*math.Pi
		}
	}
	return Lerp(a, b, t)
}
