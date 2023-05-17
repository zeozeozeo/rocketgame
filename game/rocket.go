package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

var rocketIdle = LoadEbitenImage(assets.RocketIdle)
var rocketFire = LoadEbitenImage(assets.RocketFire)
var rocketSize = Vec2i{rocketIdle.Bounds().Dx(), rocketIdle.Bounds().Dy()}

const MAX_OFFSET = 150.0
const MAX_VEL = 1.5
const DEATH_ANIM_START = 5.0
const DEATH_TIME = 7.0

type Rocket struct {
	pos               Vec2
	playerOffset      Vec2
	targetAngle       float64
	angle             float64
	vel               Vec2
	aliveTime         float64
	deathAnimProgress float64
	IsDead            bool
	didEnterBounds    bool
}

func NewRocket(x, y float64) *Rocket {
	return &Rocket{
		playerOffset: Vec2{
			RandFloat64(-MAX_OFFSET, MAX_OFFSET),
			RandFloat64(-MAX_OFFSET, MAX_OFFSET),
		},
		pos: Vec2{x, y},
	}
}

func (r *Rocket) Update(cam *Camera, player *Player, dt float64, bounds Rect) {
	r.aliveTime += dt
	if r.aliveTime > DEATH_ANIM_START {
		r.deathAnimProgress = Lerp(DEATH_ANIM_START, DEATH_TIME, r.deathAnimProgress)
	}
	r.vel.X += 0.005
	r.vel.Y += 0.005
	r.vel = r.vel.ClampMax(MAX_VEL)

	// move towards player
	r.targetAngle = RotateTowards(r.pos, player.pos.Add(r.playerOffset))
	r.angle = LerpAngle(r.angle, r.targetAngle, 0.01)
	MoveTowards(&r.pos, r.angle, r.vel)

	// check for collisions
	bounds.W += float64(rocketSize.X)
	bounds.H += float64(rocketSize.Y)
	rocketRect := r.GetRect()
	overlaps := bounds.Overlaps(rocketRect)

	if !overlaps && r.didEnterBounds {
		r.IsDead = true
	} else if overlaps {
		r.didEnterBounds = true
	}

	if rocketRect.Overlaps(player.GetRect()) {
		player.Die()
	}
}

func (r *Rocket) Draw(screen *ebiten.Image, cam *Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(rocketSize.X)/2, -float64(rocketSize.Y)/2)
	op.GeoM.Rotate(r.angle)
	op.GeoM.Translate(r.pos.X, r.pos.Y)
	cam.ApplyOP(op)

	if r.vel.X < MAX_VEL || r.vel.Y < MAX_VEL {
		screen.DrawImage(rocketFire, op)
	} else {
		screen.DrawImage(rocketIdle, op)
	}
}

func (r *Rocket) GetRect() Rect {
	return Rect{r.pos.X, r.pos.Y, float64(rocketSize.X), float64(rocketSize.Y)}
}
