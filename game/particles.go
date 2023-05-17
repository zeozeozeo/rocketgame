package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Particle struct {
	pos      Vec2
	size     Vec2
	clr      color.RGBA
	lifetime float64
}

type ParticleManager struct {
	particles []*Particle
}

func NewParticleManager() *ParticleManager {
	return &ParticleManager{}
}

func (pm *ParticleManager) Spawn(pos Vec2, size Vec2, clr color.RGBA, lifetime float64) {
	pm.particles = append(pm.particles, &Particle{pos, size, clr, lifetime})
}

func (pm *ParticleManager) Update(dt float64) {
	for i, p := range pm.particles {
		p.lifetime -= dt
		if p.lifetime <= 0.0 {
			pm.particles = append(pm.particles[:i], pm.particles[i+1:]...)
		}
	}
}

func (pm *ParticleManager) Draw(screen *ebiten.Image, cam *Camera) {
	for _, p := range pm.particles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.pos.X, p.pos.Y)
	}
}
