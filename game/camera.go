package game

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	X    float64
	Y    float64
	Zoom float64
}

func NewCamera() *Camera {
	return &Camera{
		Zoom: 5.0,
	}
}

func (cam *Camera) Update(dt float64) {

}

func (cam *Camera) MoveTo(x, y float64) {
	cam.X, cam.Y = x, y
}

func (cam *Camera) ZoomTo(zoom float64) {
	cam.Zoom = zoom
}

func (cam *Camera) ScreenToWorld(x, y float64) (float64, float64) {
	return x/cam.Zoom + cam.X/2, y/cam.Zoom + cam.Y/2
}

func (cam *Camera) GetImageOp() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cam.X/2, -cam.Y/2)
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	return op
}
