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

const (
	TURN_RATE = math.Pi / 8 // max radians per second
)

type Player struct {
	pos         Vec2
	angle       float64
	targetAngle float64
	aliveTime   float64
	vel         Vec2 // player velocity
}

func NewPlayer(spawnPos Vec2i) *Player {
	return &Player{
		pos: Vec2{X: float64(spawnPos.X), Y: float64(spawnPos.Y)},
	}
}

func (player *Player) Update(dt float64, cam *Camera) {
	player.aliveTime += dt
	if player.aliveTime < 0.6 {
		return // don't process input while respawning
	}

	player.vel.X += 0.005
	player.vel.Y += 0.005
	player.vel = player.vel.ClampMax(0.6)

	// rotate player towards mouse

	mx, my := ebiten.CursorPosition()
	wmx, wmy := cam.ScreenToWorld(float64(mx), float64(my))
	player.targetAngle = math.Atan2(float64(wmy)-player.pos.Y, float64(wmx)-player.pos.X) + math.Pi/2

	player.angle = LerpAngle(player.angle, player.targetAngle, 0.02)

	player.pos.X += math.Sin(player.angle) * player.vel.X
	player.pos.Y -= math.Cos(player.angle) * player.vel.Y
	cam.MoveTo(player.pos.X, player.pos.Y)
}

func (player *Player) Draw(screen *ebiten.Image, cam *Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(playerSize.X)/2, -float64(playerSize.Y)/2)
	op.GeoM.Rotate(player.angle)
	op.GeoM.Translate(player.pos.X, player.pos.Y)

	// make the player flicker when respawning
	if !(player.aliveTime < 0.6 && math.Mod(player.aliveTime, 0.1) < 0.05) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			screen.DrawImage(rocketFire, op)
		} else {
			screen.DrawImage(rocketIdle, op)
		}
	}
}
