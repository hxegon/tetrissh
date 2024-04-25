package tetris

import (
	"math"
	"math/rand"
)

type Piece struct {
	shape     []Vector
	color     int
	canRotate bool
}

type Color int

const (
	ColorEmpty Color = iota
	ColorGreen
	ColorYellow
	ColorRed
	ColorPurple
	ColorOrange
	ColorBlue
)

// Color is 0 for empty, or one of the 256 terminal colors
// Shouldn't this be a uint8 then? nah cause I might also make this
// a hex color code or change how this is implemented more later so nahhhh
var Pieces = []Piece{
	{ // Block
		color:     int(ColorYellow),
		canRotate: false,
		shape:     []Vector{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	},
	{ // Line
		color:     int(ColorBlue),
		canRotate: true,
		shape:     []Vector{{0, 0}, {0, -1}, {0, 1}, {0, 2}},
	},
	{ // L
		color:     int(ColorGreen),
		canRotate: true,
		shape:     []Vector{{0, -1}, {0, 0}, {0, 1}, {-1, 1}},
	},
	{ // J
		color:     int(ColorBlue),
		canRotate: true,
		shape:     []Vector{{0, -1}, {0, 0}, {0, 1}, {1, 1}},
	},
	{ // T
		color:     int(ColorPurple),
		canRotate: true,
		shape:     []Vector{{0, 0}, {-1, 0}, {1, 0}, {0, 1}},
	},
	{ // Z
		color:     int(ColorRed),
		canRotate: true,
		shape:     []Vector{{0, 0}, {-1, 0}, {0, 1}, {1, 1}},
	},
	{ // S
		color:     int(ColorOrange),
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

func (p Piece) Copy() Piece {
	newP := Piece{
		color:     p.color,
		canRotate: p.canRotate,
		shape:     make([]Vector, len(p.shape)),
	}

	copy(newP.shape, p.shape)

	return newP
}

func (p Piece) Rotate() Piece {
	newP := p.Copy()

	// idk trig I just copy paste grug
	ang := math.Pi / 2
	cos := int(math.Round(math.Cos(ang)))
	sin := int(math.Round(math.Sin(ang)))

	for i, block := range p.shape {
		nx := block.y*sin - block.x*cos
		ny := block.y*cos - block.x*sin

		newP.shape[i] = Vector{nx, ny}
	}

	return newP
}
