package game

import "golang.org/x/exp/rand"

type Vec2 struct {
	X, Y float64
}

func (v Vec2) ClampMax(max float64) Vec2 {
	if v.X > max {
		v.X = max
	}
	if v.Y > max {
		v.Y = max
	}
	return v
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}

type Vec2i struct {
	X, Y int
}

type Rect struct {
	X, Y, W, H float64
}

func (r1 Rect) Overlaps(r2 Rect) bool {
	return r1.X < r2.X+r2.W && r1.X+r1.W > r2.X && r1.Y < r2.Y+r2.H && r1.Y+r1.H > r2.Y
}

func (r Rect) SpawnRandomSide(w, h float64) Vec2 {
	v := rand.Intn(6)
	//      1
	// 0 =-----= 2
	//   -     -
	// 3 =-----= 5
	//      4
	switch v {
	case 0:
		return Vec2{r.X - w, r.Y - h}
	case 1:
		return Vec2{(r.X + r.W/2) - w/2, r.Y - h}
	case 2:
		return Vec2{r.X + r.W + w, r.Y - h}
	case 3:
		return Vec2{r.X - w, r.Y + r.H + h}
	case 4:
		return Vec2{(r.X + r.W/2) - w/2, r.Y + r.H + h}
	case 5:
		return Vec2{r.X + r.W + w, r.Y + r.H + h}
	default:
		return Vec2{}
	}
}

type Line struct {
	Start Vec2
	End   Vec2
}
