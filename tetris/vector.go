package tetris

type Vector struct {
	x, y int
}

func (v Vector) add(v2 Vector) Vector {
	return Vector{
		x: v.x + v2.x,
		y: v.y + v2.y,
	}
}
