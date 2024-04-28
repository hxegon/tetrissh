package tetris

import "fmt"

/* The tetris game should work basically like a state machine. */

type Game struct {
	board    [][]int
	piece    Piece
	pos      Vector
	height   int
	width    int
	Score    int
	GameOver bool
}

func makeBoard(height, width int) [][]int {
	board := make([][]int, height)
	blocks := make([]int, height*width)

	for i := range board {
		// Instead of allocating a new slice height times, do it all with one allocation
		// by making one []int for all blocks and chopping it up
		// This function is called ever time we want to render the board,
		// so this reduces the allocations from height+1 to 2
		board[i], blocks = blocks[:width], blocks[width:]
	}
	return board
}

// Returns false if there wasn't room for another piece
func (g *Game) NextPieceIfPossible() bool {
	piece := RandomPiece()
	pos := Vector{x: int(g.width / 2), y: 0 - piece.YOffset()}

	for _, b := range piece.shape {
		if c, ok := g.colorAt(pos.Add(b)); !ok || c > 0 {
			return false
		}
	}

	g.pos = pos
	g.piece = piece
	return true
}

func NewGame(height, width int, piece Piece) Game {
	// Initialize board
	board := makeBoard(height, width)

	g := Game{
		height:   height,
		width:    width,
		board:    board,
		GameOver: false,
	}

	g.NextPieceIfPossible()
	return g
}

func (g Game) isInBounds(v Vector) bool {
	if v.x >= 0 && v.x < g.width {
		if v.y >= 0 && v.y < g.height {
			return true
		}
	}

	return false
}

// returns the color of the block, and an OK bool for if the Vector
// was in bounds. returns 0, false for any block OOB
func (g Game) colorAt(v Vector) (int, bool) {
	if g.isInBounds(v) {
		return g.board[v.y][v.x], true
	}

	return 0, false
}

// Returns a [][]int of the board with the piece "colored" in
func (g Game) Board() [][]int {
	// make copy of board
	b := makeBoard(g.height, g.width)
	for i := range b {
		copy(b[i], g.board[i])
	}

	// set piece's shape colors on to board
	for _, v := range g.piece.shape {
		pos := g.pos.Add(v)

		// Should never happen because bounds checks exist everywhere piece is moved/placed
		if !g.isInBounds(pos) {
			msg := fmt.Sprintf("Tried to GetBoard() with a piece shape block that's out of bounds. Pos: %v, Shape: %v", g.pos, g.piece.shape)
			panic(msg)
		}

		b[pos.y][pos.x] = g.piece.color
	}

	return b
}

// Would any new positions in the shape blocks become out of bounds?
func (g *Game) moveIfPossible(direction Vector) bool {
	for _, s := range g.piece.shape {
		newPos := g.pos.Add(s).Add(direction)
		if c, inBounds := g.colorAt(newPos); !inBounds || c > 0 {
			return false
		}
	}

	g.pos = g.pos.Add(direction)
	return true
}

// argument is true for clockwise, false for counterclockwise. Returns true if it rotated, false if it couldn't
func (g *Game) rotateIfPossible() bool {
	newPiece := g.piece.Rotate()
	for _, s := range newPiece.shape {
		if c, ok := g.colorAt(g.pos.Add(s)); !ok || c > 0 {
			return false
		}
	}

	g.piece = newPiece
	return true
}

// TODO: Make this return the number of completed lines, use this for score externally instead
func (g *Game) CompactLines() {
	var broken bool
	completedLines := 0

	for y := range g.board {
		broken = false
		for x := range g.board[y] {
			if c, _ := g.colorAt(Vector{x, y}); c == 0 {
				broken = true
				break
			}
		}

		if !broken { // Line is complete
			completedLines++

			// move above lines down
			for i := y; i > 0; i-- {
				copy(g.board[i], g.board[i-1])
			}

			// add a new line to top
			g.board[0] = make([]int, len(g.board[0]))
		}
	}

	g.Score += (2 ^ completedLines*100)
}

func (g *Game) Fall() { // Maybe this should return score as well? idk
	moved := g.moveIfPossible(Vector{0, 1})

	if !moved { // Then we've reached the bottom
		g.board = g.Board()
		g.CompactLines()
		if !g.NextPieceIfPossible() {
			g.GameOver = true
		}
	}
}

// Instantly fall
func (g *Game) Drop() {
	for {
		if !g.moveIfPossible(Vector{0, 1}) {
			g.Fall()
			break
		}
	}
}

type Action int

const (
	ActionLeft Action = iota
	ActionRight
	ActionDown
	ActionRotate
	ActionDrop
)

func (g *Game) Act(a Action) {
	switch a {
	case ActionRight:
		g.moveIfPossible(Vector{1, 0})
	case ActionLeft:
		g.moveIfPossible(Vector{-1, 0})
	case ActionDown:
		g.moveIfPossible(Vector{0, 1})
	case ActionRotate:
		g.rotateIfPossible()
	case ActionDrop:
		g.Drop()
	}
}

// func (b BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {}
// func (m BoardModel) View() string {}
