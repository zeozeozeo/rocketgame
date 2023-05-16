package assets

import (
	_ "embed"
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
