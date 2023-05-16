package game

import (
	_ "embed"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

var planeAnim = []*ebiten.Image{
	LoadEbitenImage(assets.PlaneWingLeft),
	LoadEbitenImage(assets.PlaneWingRight),
	LoadEbitenImage(assets.PlaneWingBoth),
}
var playerSize = Vec2i{X: planeAnim[0].Bounds().Dx(), Y: planeAnim[0].Bounds().Dy()}

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

func NewPlayer() *Player {
	return &Player{}
}

func (player *Player) Update(dt float64, cam *Camera) {
	player.aliveTime += dt
	if player.aliveTime < 0.6 {
		return // don't process input while respawning
	}

	player.vel.X += 0.005
	player.vel.Y += 0.005
	player.vel = player.vel.ClampMax(0.6)

	// rotate & move player towards mouse
	mx, my := ebiten.CursorPosition()
	player.targetAngle = RotateTowards(player.pos, cam.ScreenToWorld(float64(mx), float64(my)))
	player.angle = LerpAngle(player.angle, player.targetAngle, 0.07)
	MoveTowards(&player.pos, player.angle, player.vel)
	cam.MoveTo(player.pos.X, player.pos.Y)
}

func (player *Player) Draw(screen *ebiten.Image, cam *Camera) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(playerSize.X)/2, -float64(playerSize.Y)/2)
	op.GeoM.Rotate(player.angle)
	op.GeoM.Translate(player.pos.X, player.pos.Y)
	cam.ApplyOP(op)

	// make the player flicker when respawning
	if !(player.aliveTime < 0.6 && math.Mod(player.aliveTime, 0.1) < 0.05) {
		animNum := int(player.aliveTime*20.0) % 2
		if player.aliveTime < 0.6 {
			animNum = 2
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			screen.DrawImage(planeAnim[animNum], op)
		} else {
			screen.DrawImage(planeAnim[animNum], op)
		}
	}
}
