package game

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
