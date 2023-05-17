package assets

import (
	_ "embed"
	"fmt"

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

var FontLoaded bool
var FutilePro font.Face

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
	FontLoaded = true
	fmt.Println("loaded font")
}
