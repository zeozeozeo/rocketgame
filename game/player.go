package game

import (
	_ "embed"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

var rocketIdle = LoadEbitenImage(assets.RocketIdle)
var rocketFire = LoadEbitenImage(assets.RocketFire)
var playerSize = Vec2i{X: rocketIdle.Bounds().Dx(), Y: rocketIdle.Bounds().Dy()}

type Player struct {
	pos       Vec2
	angle     float64
	aliveTime float64
}

func NewPlayer() *Player {
	return &Player{pos: Vec2{X: 100, Y: 100}}
}

func (player *Player) Update(dt float64, cam *Camera) {
	player.aliveTime += dt
	if player.aliveTime < 0.6 {
		return // don't process input while respawning
	}

	// rotate player towards mouse
	mx, my := ebiten.CursorPosition()
	wmx, wmy := cam.ScreenToWorld(float64(mx), float64(my))
	player.angle = math.Atan2(float64(wmy)-player.pos.Y, float64(wmx)-player.pos.X) + math.Pi/2
}

func (player *Player) Draw(screen *ebiten.Image, cam *Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(playerSize.X)/2, -float64(playerSize.Y)/2)
	op.GeoM.Rotate(player.angle)
	op.GeoM.Translate(player.pos.X, player.pos.Y)

	// make the player flicker when respawning
	if !(player.aliveTime < 0.6 && math.Mod(player.aliveTime, 0.1) < 0.05) {
		screen.DrawImage(rocketIdle, op)
	}
}
