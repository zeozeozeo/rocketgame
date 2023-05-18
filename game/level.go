package game

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zeozeozeo/rocketgame/game/assets"
)

const (
	MAX_ROCKETS = 5
	DEBUG       = false
)

var cloudTex = LoadEbitenImage(assets.Cloud)
var cloudSize = Vec2i{cloudTex.Bounds().Dx(), cloudTex.Bounds().Dy()}

type Cloud struct {
	pos            Vec2
	didEnterBounds bool
}

func (c *Cloud) GetRect() Rect {
	return Rect{c.pos.X, c.pos.Y, float64(cloudSize.X), float64(cloudSize.Y)}
}

type Level struct {
	cam            *Camera // camera
	prevCamPos     Vec2
	time           float64 // time spent in level, in seconds
	lastRocketTime float64
	player         *Player // player
	rockets        []*Rocket
	clouds         []*Cloud
	lastScoreTime  float64
	lastCloudTime  float64
	Score          int
	bestScore      int
	pm             *ParticleManager
}

func NewLevel(bestScore int) *Level {
	return &Level{
		cam:       NewCamera(),
		player:    NewPlayer(),
		pm:        NewParticleManager(),
		Score:     1,
		bestScore: bestScore,
	}
}

func (level *Level) AddRocket() {
	level.rockets = append(level.rockets, NewRocket(level.cam))
}

func (level *Level) Update(dt float64) {
	level.prevCamPos = Vec2{level.cam.X, level.cam.Y}
	level.time += dt
	level.player.Update(dt, level.cam, level.pm)
	level.cam.Zoom = 5.0

	if level.time-level.lastRocketTime > RandFloat64(1.0, 3.0) && len(level.rockets) < MAX_ROCKETS {
		level.AddRocket()
		level.lastRocketTime = level.time
	}
	if level.time-level.lastScoreTime > 0.2 {
		level.Score++
		level.lastScoreTime = level.time
	}

	// get camera bounds
	bounds := level.cam.GetBounds()

	for i := 0; i < len(level.rockets); i++ {
		r := level.rockets[i]
		r.Update(level.cam, level.player, dt, bounds, level.pm)
		if r.IsDead {
			level.rockets = append(level.rockets[:i], level.rockets[i+1:]...)
			i--
		}
	}
	level.pm.Update(dt)

	// spawn random cloud
	if level.time-level.lastCloudTime > 0.35 {
		level.SpawnCloud(bounds)
	}

	// update clouds
	for i := 0; i < len(level.clouds); i++ {
		c := level.clouds[i]
		overlaps := bounds.Overlaps(c.GetRect())

		if overlaps && !c.didEnterBounds {
			c.didEnterBounds = true
			continue
		}
		if !overlaps && c.didEnterBounds {
			level.clouds = append(level.clouds[:i], level.clouds[i+1:]...)
			i--
		}
	}
}

func (level *Level) SpawnCloud(bounds Rect) {
	cloudPos := bounds.SpawnRandomSide(float64(cloudSize.X), float64(cloudSize.Y), true)

	cloud := &Cloud{cloudPos, false}
	// make sure this cloud doesn't overlap any other clouds
	for _, c := range level.clouds {
		if c.GetRect().Overlaps(cloud.GetRect()) {
			return
		}
	}

	level.clouds = append(level.clouds, cloud)
	level.lastCloudTime = level.time
}

func (level *Level) Draw(screen *ebiten.Image) {
	cam := level.cam

	// draw clouds
	screen.Fill(color.RGBA{163, 222, 247, 255})
	for _, c := range level.clouds {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.pos.X, c.pos.Y)
		level.cam.ApplyOP(op)
		screen.DrawImage(cloudTex, op)
	}
	/*
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		level.cam.ApplyOP(op)
		screen.DrawImage(cloudTex, op)
	*/

	level.pm.Draw(screen, cam)
	level.player.Draw(screen, cam)
	for _, r := range level.rockets {
		r.Draw(screen, cam)
	}

	// draw score
	sx, _ := ebiten.WindowSize()
	scoreText := fmt.Sprintf("%d", level.Score)
	tw, ty := MeasureText(scoreText, false)
	DrawTextShadow(screen, scoreText, sx/2-tw/2, 24+ty, color.RGBA{255, 255, 255, 255}, false)

	// draw best score
	if level.bestScore != 0 {
		scoreText = fmt.Sprintf("%d", level.bestScore)
		tw, ty = MeasureText(scoreText, true)
		DrawTextShadow(screen, scoreText, sx/2-tw/2, 40+ty*2, color.RGBA{255, 255, 255, 255}, true)
	}

	if DEBUG {
		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf(
				"fps: %f\ntps: %f\n%d rockets\nscore: %d",
				ebiten.ActualFPS(),
				ebiten.ActualTPS(),
				len(level.rockets),
				level.Score,
			),
		)
	}
}

func (level *Level) IsDone() bool {
	return level.player.IsDead()
}

func (level *Level) PlayRespawnSound() {
	level.player.PlayRespawnSound()
}
