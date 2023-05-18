package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Camera struct {
	X              float64
	Y              float64
	Zoom           float64
	totalShakeTime float64
	shakeTime      float64
	lastShakeTime  float64
	isShaking      bool
	lastShakeVal   Vec2
}

func NewCamera() *Camera {
	return &Camera{
		Zoom: 1.0,
	}
}

func (cam *Camera) Shake(time float64) {
	cam.totalShakeTime = time
	cam.shakeTime = time
}

func (cam *Camera) Update(dt float64) {
	if cam.shakeTime > 0.0 {
		cam.shakeTime -= dt
	}
	cam.isShaking = cam.shakeTime > 0.0
}

func (cam *Camera) MoveTo(x, y float64) *Camera {
	cam.X, cam.Y = x, y
	return cam
}

func (cam *Camera) SetZoom(zoom float64) *Camera {
	cam.Zoom = zoom
	return cam
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

func (cam *Camera) getShakeVal() (shakeX, shakeY float64) {
	if cam.isShaking && math.Abs(cam.lastShakeTime-cam.shakeTime) > 0.05 {
		shakeProg := 1.0 - Normalize(cam.shakeTime, cam.totalShakeTime, 0.0) + 0.2
		shakeX = (RandFloat64(-5.0, 5.0) / cam.Zoom) / shakeProg
		shakeY = (RandFloat64(-5.0, 5.0) / cam.Zoom) / shakeProg
		cam.lastShakeVal = Vec2{shakeX, shakeY}
		cam.lastShakeTime = cam.shakeTime
	} else {
		shakeX = cam.lastShakeVal.X
		shakeY = cam.lastShakeVal.Y
	}
	return
}

func (cam *Camera) ApplyOP(op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	wx, wy := ebiten.WindowSize()
	div := cam.Zoom * 2
	shakeX, shakeY := cam.getShakeVal()
	op.GeoM.Translate(-cam.X+(float64(wx)/div)+shakeX, -cam.Y+(float64(wy)/div)+shakeY)
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	return op
}

func (cam *Camera) ApplyOPColorM(op *colorm.DrawImageOptions) *colorm.DrawImageOptions {
	wx, wy := ebiten.WindowSize()
	div := cam.Zoom * 2
	shakeX, shakeY := cam.getShakeVal()
	op.GeoM.Translate(-cam.X+(float64(wx)/div)+shakeX, -cam.Y+(float64(wy)/div)+shakeY)
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	return op
}

func (cam *Camera) GetImageOp() *ebiten.DrawImageOptions {
	return cam.ApplyOP(&ebiten.DrawImageOptions{})
}
