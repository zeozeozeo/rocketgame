package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Camera struct {
	X    float64
	Y    float64
	Zoom float64
}

func NewCamera() *Camera {
	return &Camera{
		Zoom: 1.0,
	}
}

func (cam *Camera) Update(dt float64) {

}

func (cam *Camera) MoveTo(x, y float64) {
	cam.X, cam.Y = x, y
}

func (cam *Camera) SetZoom(zoom float64) {
	cam.Zoom = zoom
}

func (cam *Camera) ScreenToWorld(x, y float64) Vec2 {
	wx, wy := ebiten.WindowSize()
	div := cam.Zoom * 2
	return Vec2{
		x/cam.Zoom + cam.X - (float64(wx) / div),
		y/cam.Zoom + cam.Y - (float64(wy) / div),
	}
}

func (cam *Camera) WorldToScreen(x, y float64) Vec2 {
	wx, wy := ebiten.WindowSize()
	return Vec2{
		(x-cam.X)*cam.Zoom + (float64(wx) / 2),
		(y-cam.Y)*cam.Zoom + (float64(wy) / 2),
	}
}

func (cam *Camera) GetBounds() Rect {
	sx, sy := ebiten.WindowSize()
	start := cam.ScreenToWorld(0, 0)
	end := cam.ScreenToWorld(float64(sx), float64(sy))
	return Rect{
		start.X,
		start.Y,
		end.X - start.X,
		end.Y - start.Y,
	}
}

func (cam *Camera) ApplyOP(op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	wx, wy := ebiten.WindowSize()
	div := cam.Zoom * 2
	op.GeoM.Translate(-cam.X+(float64(wx)/div), -cam.Y+(float64(wy)/div))
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	return op
}

func (cam *Camera) ApplyOPColorM(op *colorm.DrawImageOptions) *colorm.DrawImageOptions {
	wx, wy := ebiten.WindowSize()
	div := cam.Zoom * 2
	op.GeoM.Translate(-cam.X+(float64(wx)/div), -cam.Y+(float64(wy)/div))
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	return op
}

func (cam *Camera) GetImageOp() *ebiten.DrawImageOptions {
	return cam.ApplyOP(&ebiten.DrawImageOptions{})
}
