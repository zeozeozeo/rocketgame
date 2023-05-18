package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/zeozeozeo/rocketgame/game"
)

const TPS = 144

type Game struct {
	level        *game.Level
	bestScore    int
	isFullscreen bool
	prevSize     game.Vec2i
}

func (g *Game) Update() error {
	// restart game if the level is done
	if g.level.IsDone() {
		score := g.level.Score
		if score > g.bestScore {
			g.bestScore = score
		}

		g.level = game.NewLevel(g.bestScore)
		g.level.PlayRespawnSound()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.isFullscreen = !g.isFullscreen
		if !g.isFullscreen {
			ebiten.SetWindowSize(g.prevSize.X, g.prevSize.Y)
		} else {
			wx, wy := ebiten.WindowSize()
			g.prevSize = game.Vec2i{X: wx, Y: wy}
		}
		ebiten.SetFullscreen(g.isFullscreen)
	}

	g.level.Update(1.0 / TPS)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ebiten.SetWindowSize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(TPS)
	ebiten.SetWindowTitle("rocketgame")

	g := &Game{}
	g.level = game.NewLevel(0)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
