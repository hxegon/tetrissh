package tetris

import "math/rand"

type Piece struct {
	shape     []Vector
	color     int
	canRotate bool
}

var Pieces = []Piece{
	{ // Block
		color:     1,
		canRotate: false,
		shape:     []Vector{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	},
	{ // Line
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {0, -1}, {0, 1}, {0, 2}},
	},
	{ // L
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {0, 1}, {0, 2}, {-1, 2}},
	},
	{ // J
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
	},
	{ // T
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {-1, 0}, {1, 0}, {0, 1}},
	},
	{ // Z
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {-1, 0}, {0, 1}, {1, 1}},
	},
	{ // S
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {1, 0}, {0, 1}, {-1, 1}},
	},
}

func RandomPiece() Piece {
	idx := rand.Intn(len(Pieces))
	return Pieces[idx]
}

func (p Piece) YOffset() int {
	offset := 0

	for _, v := range p.shape {
		if v.y < offset {
			offset = v.y
		}
	}

	return offset
}

// func NewRandomPiece() Piece {}
// func (p *Piece) Rotate() {}
// func (p *Piece) RotateBack() {}
