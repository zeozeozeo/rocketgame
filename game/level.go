package game

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	lastScoreTime  float64
	score          int
	pm             *ParticleManager
}

func NewLevel() *Level {
	return &Level{
		cam:    NewCamera(),
		player: NewPlayer(),
		pm:     NewParticleManager(),
		score:  1,
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
	if level.time-level.lastScoreTime > 0.2 {
		level.score++
		level.lastScoreTime = level.time
	}

	// get camera bounds
	bounds := level.cam.GetBounds()

	for i, r := range level.rockets {
		r.Update(level.cam, level.player, dt, bounds)
		if r.IsDead {
			level.rockets = append(level.rockets[:i], level.rockets[i+1:]...)
		}
	}
}

func (level *Level) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{163, 222, 247, 255})

	level.player.Draw(screen, level.cam)
	for _, r := range level.rockets {
		r.Draw(screen, level.cam)
	}

	// draw score
	scoreText := fmt.Sprintf("%d", level.score)
	sx, _ := ebiten.WindowSize()
	tw, ty := MeasureText(scoreText)
	DrawTextShadow(screen, scoreText, sx/2-tw/2, 24+ty, color.RGBA{255, 255, 255, 255})

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"fps: %f\ntps: %f\n%d rockets\nscore: %d",
			ebiten.ActualFPS(),
			ebiten.ActualTPS(),
			len(level.rockets),
			level.score,
		),
	)
}

func (level *Level) IsDone() bool {
	return level.player.IsDead()
}
