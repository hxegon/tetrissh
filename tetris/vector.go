package tetris

type Vector struct {
	x, y int
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		x: v.x + v2.x,
		y: v.y + v2.y,
	}
}
