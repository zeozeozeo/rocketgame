package game

import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Block uint8

const (
	BLOCK_AIR   Block = iota
	BLOCK_SOLID Block = iota
)

type Level struct {
	cam            *Camera // camera
	time           float64 // time spent in level, in seconds
	lastRocketTime float64
	player         *Player // player
	rockets        []*Rocket
}

func NewLevel() *Level {
	return &Level{
		cam:    NewCamera(),
		player: NewPlayer(),
	}
}

func (level *Level) AddRocket() {
	level.rockets = append(level.rockets, NewRocket(
		level.player.pos.X+RandFloat64(-100.0, 100.0),
		level.player.pos.Y+100.0,
	))
}

func (level *Level) Update(dt float64) {
	level.time += dt
	level.player.Update(dt, level.cam)
	level.cam.Zoom = 5.0

	if level.time-level.lastRocketTime > 1.0 && len(level.rockets) < 5 {
		level.AddRocket()
		level.lastRocketTime = level.time
	}

	for _, r := range level.rockets {
		r.Update(level.cam, level.player, dt)
	}
}

func (level *Level) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{163, 222, 247, 255})

	for _, r := range level.rockets {
		r.Draw(screen, level.cam)
	}

	level.player.Draw(screen, level.cam)
}
