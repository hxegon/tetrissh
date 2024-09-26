package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type matchState int

const (
	msLooking matchState = iota
	msRunning
	msCanceled
)

func (m matchState) String() string {
	switch m {
	case msLooking:
		return "looking for match"
	case msRunning:
		return "match active"
	case msCanceled:
		return "match canceled"
	default:
		return "invalid matchState"
	}
}

type MultiplayerGame struct {
	cancel    context.CancelFunc
	session   *MultiplayerSession
	opSession *MultiplayerSession
	opC       <-chan *MultiplayerSession // Channel for the matchmaking goroutine to send the opponent to
	game      *GameModel
	scoreBar  *progress.Model
	mstate    matchState
}

func NewMultiplayer() *MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int

	// TODO: NewMultiplayerSession
	session := &MultiplayerSession{
		ctx:   ctx,
		board: board,
		mx:    new(sync.RWMutex),
	}

	gm := NewGameModel()
	// initialize the board pointer, it shouldn't be nil unless the game has been closed
	session.SetBoard(gm.Board())

	opC := session.requestMatch()

	scoreM := progress.New(
		progress.WithGradient("#030ffc", "#fa0202"),
		progress.WithoutPercentage(),
	)

	game := &MultiplayerGame{
		session:  session,
		cancel:   cancel,
		scoreBar: &scoreM,
		opC:      opC,
		game:     &gm,
	}

	return game
}

func (m *MultiplayerGame) close() {
	log.Debug("Closing game")
	m.cancel()
	// drop shared pointers, might not be necessary
	m.opSession = nil
	m.session = nil
}

// Return enum checking if the game is looking for a match, running or canceled
// sets and returns matchState, but m.mstate shouldn't be accessed other than through here
func (m *MultiplayerGame) setState() {
	oldState := m.mstate

	if oldState != msCanceled {
		newState := m.mstate

		if m.opSession == nil {
			// Why does this throw nil pointer deref error if oldState == msRunning is inlined with above?
			if oldState == msRunning {
				newState = msCanceled
			}
		} else {
			select {
			// Check if either session in the match is canceled
			case <-m.opSession.done():
				newState = msCanceled
			case <-m.session.done(): // Shouldn't really be reached but just in case
				newState = msCanceled
			default: // If not, make sure mstate is running
				if oldState == msLooking {
					log.Debug("Setting multiplayer game to running, initializing fall tick")
					newState = msRunning
				}
			}
		}

		m.mstate = newState
	}
}

// Update this game's session board. Thread safe, blocks on mutex
func (m *MultiplayerGame) SetBoard() {
	m.session.SetBoard(m.game.Board())
}

type MatchLookTickMsg struct{}

func MatchLookTick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return MatchLookTickMsg{}
	})
}

func (m MultiplayerGame) Init() tea.Cmd {
	return MatchLookTick()
}

// FIXME: Doesn't properly cancel multiplayer match if session is terminated OOB (not C-c or q)
func (m MultiplayerGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.setState()

	switch msg := msg.(type) {
	case MatchLookTickMsg:
		// Continue refreshing if there's no match found
		if m.mstate == msLooking {
			// Once an opponent session has been sent, set the opponent
			select {
			case s, ok := <-m.opC:
				if !ok {
					log.Error("Tried to read from a closed opC on a MatchLook msg")
					// TODO: Return error msg cmd here
					return m, nil
				} else {
					m.opSession = s
					return m, m.game.Init()
				}
			default:
				return m, MatchLookTick()
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.close()
			return m, DeactivateCmd
		default:
			if m.mstate == msRunning {
				*m.game, cmd = m.game.Update(msg)
				m.SetBoard()
			}
		}
	case FallMsg:
		*m.game, cmd = m.game.Update(msg)
		m.SetBoard()
	}
	return m, cmd
}

// Returns a # between 0.0 <-> 1.0 representing the ration between this player's score and the other
// Returns 0.5, err if we can't get score.
func (m *MultiplayerGame) scoreRatio() (float64, error) {
	if m.session == nil || m.opSession == nil {
		msg := "Tried to call scorePercent on a game with a nil session or opSession pointer"
		return 0.5, fmt.Errorf(msg)
	}

	score := m.session.Score()
	opScore := m.opSession.Score()

	if score+opScore == 0 { // Ratio is "even" if both scores are 0. Guard against Div by 0
		return 0.5, nil
	}

	return float64(score) / float64(score+opScore), nil
}

func (m *MultiplayerGame) renderGame() string {
	// Layout is as such:
	// Score bar
	// Your score | Their score
	// Your Board | Their board
	boardsView := lipgloss.JoinHorizontal(
		lipgloss.Left,
		BoardView(m.game),
		BoardView(m.opSession),
	)

	boardW := lipgloss.Width(boardsView)
	m.scoreBar.Width = boardW

	if r, err := m.scoreRatio(); err != nil {
		panic(err)
	} else {
		scoreView := m.scoreBar.ViewAs(r)
		return lipgloss.JoinVertical(lipgloss.Left, scoreView, boardsView)
	}
}

func (m MultiplayerGame) View() string {
	m.setState()
	switch m.mstate {
	case msLooking:
		return "looking for match"
	case msRunning:
		if err := m.opSession.err; err != nil {
			// TODO: Render error message
			msg := fmt.Sprintf("Error when trying to view an opSession: %v", err)
			log.Errorf(msg)
			panic(msg) // TODO: just display the message
		}
		return m.renderGame()

	case msCanceled:
		return "match canceled"
	default:
		log.Warn("invalid match state in multiplayergame view function")
		return "something went wrong!"
	}
}
