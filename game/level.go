package game

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed levels/*
var levels embed.FS

type Block uint8

const (
	BLOCK_AIR   Block = iota
	BLOCK_SOLID Block = iota
)

type Level struct {
	grid [][]Block
	cam  *Camera
}

func (level *Level) Set(x, y int, block Block) {
	level.grid[y][x] = block
}

func (level *Level) At(x, y int) Block {
	return level.grid[y][x]
}

func LoadLevel(num int) (*Level, error) {
	start := time.Now()
	level := &Level{}

	f, err := levels.Open(fmt.Sprintf("levels/%d.png", num))
	if err != nil {
		return nil, err
	}

	// decode image
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		return nil, err
	}

	// allocate grid (by default, all blocks are air)
	level.grid = make([][]Block, img.Bounds().Dy())
	for i := range level.grid {
		level.grid[i] = make([]Block, img.Bounds().Dx())
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 0,0,0 = solid
			if r == 0 && g == 0 && b == 0 {
				level.Set(x, y, BLOCK_SOLID)
			}
		}
	}

	// create camera
	level.cam = NewCamera()

	fmt.Printf("loaded level %d in %s\n", num, time.Since(start))
	return level, nil
}

func (level *Level) Update() {

}

func (level *Level) drawBlock(screen *ebiten.Image, x, y int) {
	block := level.At(x, y)
	rx := float32(x)
	ry := float32(y)
	switch block {
	case BLOCK_SOLID:
		vector.DrawFilledRect(screen, rx, ry, 1, 1, color.RGBA{255, 255, 255, 255}, false)
	}
}

func (level *Level) Draw(screen *ebiten.Image) {
	for y := 0; y < len(level.grid); y++ {
		for x := 0; x < len(level.grid[y]); x++ {
			level.drawBlock(screen, x, y)
		}
	}
}
