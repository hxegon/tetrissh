package tetris

import (
	"testing"
)

var testPiece = Piece{shape: []Vector{{0, 0}}, color: 1, canRotate: false}

func TestIsInBounds(t *testing.T) {
	g := NewGame(10, 10, testPiece)
	cases := []struct {
		name string
		vec  Vector
		in   bool
	}{
		{"In Bounds", Vector{0, 0}, true},
		{"In Bounds", Vector{5, 5}, true},
		{"In Bounds", Vector{9, 9}, true},
		{"Too high", Vector{0, -1}, false},
		{"Too left", Vector{-1, 5}, false},
		{"Too right", Vector{11, 9}, false},
		{"Too low", Vector{0, 11}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if g.isInBounds(c.vec) != c.in {
				t.Errorf("%v's in bounds status should have been %v but wasn't", c.vec, c.in)
			}
		})
	}
}

func TestGetBoard(t *testing.T) {
	width := 10
	height := 15
	g := NewGame(height, width, testPiece)
	b := g.GetBoard()

	expectedPiecePos := Vector{x: int(width / 2), y: 0}

	for y := range b {
		for x := range b[y] {
			if !g.isInBounds(Vector{x, y}) {
				t.Fatalf("Somehow tried to test a color that's out of bounds: x: %v, y: %v", x, y)
			}

			color := b[y][x]

			// if we're looking at the block where the piece should be
			if x == expectedPiecePos.x && y == expectedPiecePos.y {
				if color != 1 { // and it's not filled in
					t.Errorf("Piece block wasn't colored")
				}
			} else { // otherwise it should be empty
				if color != 0 {
					t.Errorf("Block was colored when it should be empty: x: %v, y: %v", x, y)
				}
			}
		}
	}
}
