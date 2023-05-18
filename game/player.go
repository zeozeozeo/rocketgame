package game

import (
	_ "embed"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	pos             Vec2
	angle           float64
	targetAngle     float64
	aliveTime       float64
	vel             Vec2 // player velocity
	isDead          bool
	trail           [256]Line
	curTrail        int
	prevTrailPos    Vec2
	frame           uint64
	deathAnimActive bool
	deathAnimTimer  float64
}

func NewPlayer() *Player {
	return &Player{deathAnimTimer: EXPLOSION_ANIM_LENGTH}
}

func (player *Player) Update(dt float64, cam *Camera, pm *ParticleManager) {
	if player.deathAnimTimer <= 0 {
		player.isDead = true
		return
	}
	if player.deathAnimActive {
		player.deathAnimTimer -= dt
		return
	}

	player.frame++
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

	// process trail
	if player.frame%5 == 0 {
		player.trail[player.curTrail] = Line{
			player.prevTrailPos,
			player.pos,
		}
		player.curTrail++
		player.curTrail %= len(player.trail)
		player.prevTrailPos = player.pos
	}
}

func (player *Player) Draw(screen *ebiten.Image, cam *Camera) {
	if player.deathAnimActive {
		return
	}

	// draw trail
	for i := 0; i < len(player.trail); i++ {
		if i%3 != 0 {
			continue
		}

		start := player.trail[i].Start
		end := player.trail[i].End
		startW := cam.WorldToScreen(start.X, start.Y)
		endW := cam.WorldToScreen(end.X, end.Y)
		vector.StrokeLine(
			screen,
			float32(startW.X),
			float32(startW.Y),
			float32(endW.X),
			float32(endW.Y),
			float32(cam.Zoom),
			color.RGBA{255, 255, 255, 255},
			false,
		)
	}

	// draw player
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

func (player *Player) GetRect() Rect {
	width, height := float64(playerSize.X)/2, float64(playerSize.Y)/2
	return Rect{player.pos.X + width/2, player.pos.Y + height/2, width, height}
}

func (player *Player) Die(cam *Camera, pm *ParticleManager) {
	if player.deathAnimActive {
		return
	}
	NewSound(assets.ExplosionSound).SetVolume(0.06).Play()

	cam.Shake(EXPLOSION_ANIM_LENGTH)
	pm.SpawnExplosion(player.pos)
	player.deathAnimActive = true
}

func (player *Player) IsDead() bool {
	return player.isDead
}

func (player *Player) PlayRespawnSound() {
	NewSound(assets.RespawnSound).SetVolume(0.1).Play()
}
