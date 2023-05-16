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

// func (v1 Vec2) Sub(v2 Vec2) Vec2 {
// 	  return Vec2{v1.X - v2.X, v1.Y - v2.Y}
// }

type Vec2i struct {
	X, Y int
}
