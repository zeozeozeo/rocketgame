package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zeozeozeo/rocketgame/game"
)

const WIDTH, HEIGHT = 1280, 720
const TPS = 144

type Game struct {
	level *game.Level
}

func (g *Game) Update() error {
	// restart game if the level is done
	if g.level.IsDone() {
		g.level = game.NewLevel()
	}

	g.level.Update(1.0 / TPS)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.Draw(screen)
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
