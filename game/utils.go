package game

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"io"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/zeozeozeo/rocketgame/game/assets"
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

func shortAngleDist(a float64, b float64) float64 {
	turn := math.Pi * 2
	deltaAngle := math.Mod(b-a, turn)
	return math.Mod(2*deltaAngle, turn) - deltaAngle
}

func LerpAngle(a, b, t float64) float64 {
	return a + shortAngleDist(a, b)*t
}

func RotateTowards(from, to Vec2) float64 {
	return math.Atan2(float64(to.Y)-from.Y, float64(to.X)-from.X) + math.Pi/2
}

func MoveTowards(pos *Vec2, angle float64, vel Vec2) {
	pos.X += math.Sin(angle) * vel.X
	pos.Y -= math.Cos(angle) * vel.Y
}

func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func ClampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func Normalize(value, max, min float64) float64 {
	return (value - min) / (max - min)
}

func DrawText(screen *ebiten.Image, str string, x, y int, clr color.RGBA, small bool) {
	if !assets.FontLoaded {
		assets.LoadFont()
	}
	if small {
		text.Draw(screen, str, assets.FutileProSmall, x, y, clr)
	} else {
		text.Draw(screen, str, assets.FutilePro, x, y, clr)
	}
}

func MeasureText(str string, small bool) (int, int) {
	if !assets.FontLoaded {
		assets.LoadFont()
	}

	var rect image.Rectangle
	if small {
		rect = text.BoundString(assets.FutileProSmall, str)
	} else {
		rect = text.BoundString(assets.FutilePro, str)
	}

	return rect.Dx(), rect.Dy()
}

func DrawTextShadow(screen *ebiten.Image, str string, x, y int, clr color.RGBA, small bool) {
	DrawText(screen, str, x+2, y+2, color.RGBA{0, 0, 0, 255}, small)
	DrawText(screen, str, x, y, clr, small)
}

var fireColors = []color.RGBA{
	{255, 242, 0, 255},
	{252, 182, 6, 255},
	{236, 28, 36, 255},
	{236, 194, 14, 255},
	{255, 230, 0, 255},
	{255, 214, 0, 255},
	{255, 198, 0, 255},
	{255, 178, 0, 255},
}

func RandomFireColor() color.RGBA {
	return fireColors[rand.Intn(len(fireColors))]
}
