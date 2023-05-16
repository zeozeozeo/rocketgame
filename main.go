package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zeozeozeo/rocketgame/game"
)

const WIDTH, HEIGHT = 1280, 720
const TPS = 144

type Game struct {
	level *game.Level
}

func (g *Game) Update() error {
	g.level.Update(1.0 / TPS)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %f\ntps: %f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetTPS(TPS)
	ebiten.SetWindowTitle("rocketgame")

	g := &Game{}
	g.level = game.NewLevel()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
