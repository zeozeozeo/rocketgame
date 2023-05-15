package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zeozeozeo/rocketgame/game"
)

const width, height = 1280, 720

type Game struct {
	currentLevel *game.Level
}

func (g *Game) Update() error {
	g.currentLevel.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentLevel.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowTitle("rocketgame")

	g := &Game{}
	level, err := game.LoadLevel(1)
	if err != nil {
		panic(err)
	}
	g.currentLevel = level

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
