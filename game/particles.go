package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

var particleTex = LoadEbitenImage(assets.Particle)
var particleSize = Vec2i{particleTex.Bounds().Dx(), particleTex.Bounds().Dy()}

const DECAY_TIME = 0.2

type Particle struct {
	pos    Vec2
	size   Vec2
	vel    Vec2
	rotVel float64
	rot    float64
	clr    color.RGBA
	life   float64
}

type ParticleManager struct {
	particles []*Particle
}

func NewParticleManager() *ParticleManager {
	return &ParticleManager{}
}

func (pm *ParticleManager) Spawn(pos Vec2, size Vec2, vel Vec2, rotVel float64, clr color.RGBA, lifetime float64) {
	pm.particles = append(pm.particles, &Particle{pos, size, vel, rotVel, 0.0, clr, lifetime})
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
}
