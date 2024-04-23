package tetris

type Piece struct {
	shape     []Vector
	color     int
	canRotate bool
}

var Pieces = []Piece{
	{ // unit block for testing
		color:     1,
		canRotate: false,
		shape:     []Vector{{0, 0}},
	},
	{ // Line
		color:     1,
		canRotate: true,
		shape:     []Vector{{0, 0}, {0, -1}, {0, 1}, {0, 2}},
	},
}

// func NewRandomPiece() Piece {}
// func (p *Piece) YOffset() int {}
// func (p *Piece) Rotate() {}
// func (p *Piece) RotateBack() {}
