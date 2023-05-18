package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

var particleTex = LoadEbitenImage(assets.Particle)
var particleSize = Vec2i{particleTex.Bounds().Dx(), particleTex.Bounds().Dy()}
var explosionTex = LoadEbitenImage(assets.Explosion)

const (
	DECAY_TIME            = 0.2
	EXPLOSION_FRAME_TIME  = 0.03
	EXPLOSION_FRAMES      = 17
	EXPLOSION_SIZE        = 192
	EXPLOSION_ANIM_LENGTH = (EXPLOSION_FRAMES + 4) * EXPLOSION_FRAME_TIME
)

type Particle struct {
	pos    Vec2
	size   Vec2
	vel    Vec2
	rotVel float64
	rot    float64
	clr    color.RGBA
	life   float64
}

type Explosion struct {
	pos           Vec2
	nextFrameTime float64
	frame         int
}

type ParticleManager struct {
	particles  []*Particle
	explosions []*Explosion
}

func NewParticleManager() *ParticleManager {
	return &ParticleManager{}
}

func (pm *ParticleManager) Spawn(pos Vec2, size Vec2, vel Vec2, rotVel float64, clr color.RGBA, lifetime float64) {
	pm.particles = append(pm.particles, &Particle{pos, size, vel, rotVel, 0.0, clr, lifetime})
}

func (pm *ParticleManager) SpawnExplosion(pos Vec2) {
	pm.explosions = append(pm.explosions, &Explosion{pos: pos, nextFrameTime: EXPLOSION_FRAME_TIME})
}

func (pm *ParticleManager) Update(dt float64) {
	for i := 0; i < len(pm.particles); i++ {
		p := pm.particles[i]
		p.life -= dt
		if p.life <= 0.0 {
			pm.particles = append(pm.particles[:i], pm.particles[i+1:]...)
			i--
		}
		p.pos.X += p.vel.X
		p.pos.Y += p.vel.Y
		p.rot += p.rotVel
	}

	for i := 0; i < len(pm.explosions); i++ {
		exp := pm.explosions[i]

		if exp.frame >= EXPLOSION_FRAMES {
			pm.explosions = append(pm.explosions[:i], pm.explosions[i+1:]...)
			i--
			continue
		}

		exp.nextFrameTime -= dt
		if exp.nextFrameTime <= 0.0 {
			exp.frame++
			exp.nextFrameTime = EXPLOSION_FRAME_TIME
		}
	}
}

func (pm *ParticleManager) Draw(screen *ebiten.Image, cam *Camera) {
	for _, p := range pm.particles {
		var cm colorm.ColorM
		cm.ScaleWithColor(p.clr)

		// decrease transparency
		if p.life <= DECAY_TIME {
			cm.Scale(1.0, 1.0, 1.0, Normalize(DECAY_TIME-p.life, 0.0, DECAY_TIME))
		}

		op := &colorm.DrawImageOptions{}

		op.GeoM.Scale(p.size.X, p.size.Y)
		op.GeoM.Translate(-float64(particleSize.X)/2, -float64(particleSize.Y)/2)
		op.GeoM.Rotate(p.rot)

		op.GeoM.Translate(p.pos.X, p.pos.Y)
		cam.ApplyOPColorM(op)
		colorm.DrawImage(screen, particleTex, cm, op)
	}
	for _, exp := range pm.explosions {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(EXPLOSION_SIZE)/2, -float64(EXPLOSION_SIZE)/2)
		op.GeoM.Scale(0.2, 0.2)
		op.GeoM.Translate(exp.pos.X, exp.pos.Y)
		cam.ApplyOP(op)

		frame := exp.frame
		tex := explosionTex.SubImage(image.Rect(
			EXPLOSION_SIZE*frame,
			0,
			EXPLOSION_SIZE*(frame+1),
			EXPLOSION_SIZE,
		)).(*ebiten.Image)
		screen.DrawImage(tex, op)
	}
}
