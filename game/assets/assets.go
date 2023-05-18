package assets

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed plane_wing_left.png
var PlaneWingLeft []byte

//go:embed plane_wing_both.png
var PlaneWingBoth []byte

//go:embed plane_wing_right.png
var PlaneWingRight []byte

//go:embed rocket_idle.png
var RocketIdle []byte

//go:embed rocket_fire.png
var RocketFire []byte

//go:embed FutilePro.ttf
var futileProFont []byte

//go:embed particle.png
var Particle []byte

//go:embed cloud.png
var Cloud []byte

//go:embed explosion.png
var Explosion []byte

//go:embed explosion.wav
var ExplosionSound []byte

//go:embed respawn.wav
var RespawnSound []byte

var FontLoaded bool
var FutilePro font.Face
var FutileProSmall font.Face

func LoadFont() {
	tt, err := opentype.Parse(futileProFont)
	if err != nil {
		panic(err)
	}

	const dpi = 72
	FutilePro, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	FutileProSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	FontLoaded = true
}
