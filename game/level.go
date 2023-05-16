package game

import (
	"embed"
	"fmt"
	"image/color"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed levels/*
var levels embed.FS

type Block uint8

const (
	BLOCK_AIR   Block = iota
	BLOCK_SOLID Block = iota
)

type Level struct {
	grid       [][]Block
	levelImage *ebiten.Image // base level image, doesn't change
	mapImage   *ebiten.Image // map image, rendered every frame
	cam        *Camera       // camera
	time       float64       // time spent in level, in seconds
	spawnPos   Vec2i         // player spawn position
	player     *Player       // player
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
	img, err := LoadRGBAImage(f)
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
			// 0-65535 => 0-255
			r >>= 8
			g >>= 8
			b >>= 8
			clearPixel := false

			// black = solid
			// blue = spawn position

			if r == 0 && g == 0 && b == 0 {
				level.Set(x, y, BLOCK_SOLID)
			}
			if r == 0 && g == 0 && b == 255 {
				fmt.Println("found spawn position at", x, y)
				level.spawnPos = Vec2i{x, y}
				clearPixel = true
			}

			if clearPixel {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// create camera, level and player
	level.cam = NewCamera()
	level.levelImage = ebiten.NewImageFromImage(img)
	level.mapImage = ebiten.NewImage(img.Bounds().Dx(), img.Bounds().Dy())
	level.player = NewPlayer(level.spawnPos)

	fmt.Printf("loaded level %d in %s\n", num, time.Since(start))
	return level, nil
}

func (level *Level) Update(dt float64) {
	level.time += dt
	level.player.Update(dt, level.cam)
	level.cam.Zoom = 5.0
}

func (level *Level) drawMap() {
	level.mapImage.DrawImage(level.levelImage, &ebiten.DrawImageOptions{})
	level.player.Draw(level.mapImage, level.cam)
}

func (level *Level) drawMapToScreen(screen *ebiten.Image) {
	screen.DrawImage(level.mapImage, level.cam.GetImageOp())
}

func (level *Level) Draw(screen *ebiten.Image) {
	level.drawMap()
	level.drawMapToScreen(screen)
}
