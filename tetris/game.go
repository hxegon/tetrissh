package tetris

/* The tetris game should work basically like a state machine. */

type Game struct {
	board  [][]int
	piece  Piece
	pos    Vector
	height int
	width  int
}

func makeBoard(height, width int) [][]int {
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}
	return board
}

func (g *Game) SetRandomPiece() {
	g.piece = RandomPiece()
	g.pos = Vector{x: int(g.width / 2), y: 0 - g.piece.YOffset()}
}

func NewGame(height, width int, piece Piece) Game {
	// Initialize board
	board := makeBoard(height, width)

	g := Game{
		height: height,
		width:  width,
		board:  board,
	}

	g.SetRandomPiece()
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
func (g Game) GetBoard() [][]int {
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
			panic("Tried to get a tetris board with an active piece that's out of bounds. This should never happen!")
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

func (g *Game) Fall() { // This should be different than the user's "down" action
	// TODO: detect lines and compact
	// TODO: Score
}

type Action int

const (
	ActionLeft Action = iota
	ActionRight
	ActionDown
	ActionRotate
	ActionRotateBack
)

func (g *Game) Act(a Action) {
	switch a {
	case ActionRight:
		g.moveIfPossible(Vector{1, 0})
	case ActionLeft:
		g.moveIfPossible(Vector{-1, 0})
	}
}

// func (b BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {}
// func (m BoardModel) View() string {}
