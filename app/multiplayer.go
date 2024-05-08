package app

import (
	"context"
	"sync"
	"time"

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
	opC       <-chan *MultiplayerSession
	game      *GameModel
	mstate    matchState
}

func NewMultiplayer() *MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int

	session := &MultiplayerSession{
		ctx:   ctx,
		board: board,
		mx:    new(sync.RWMutex),
	}

	gm := NewGameModel()
	// initialize the board pointer, it shouldn't be nil unless the game has been closed
	session.SetBoard(gm.Board())

	opC := session.requestMatch()

	game := &MultiplayerGame{
		session: session,
		cancel:  cancel,
		opC:     opC,
		game:    &gm,
	}

	return game
}

func (m *MultiplayerGame) close() {
	log.Debug("Closing game")
	m.cancel()
	// drop pointers
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

func (m *MultiplayerGame) SetBoard() {
	// FIXME: is session context closed? what about opSession? etc.
	m.session.SetBoard(m.game.Board())
}

func (m *MultiplayerGame) opBoard() [][]int {
	if b, ok := m.opSession.Board(); ok {
		return b
	}
	panic("couldn't get board")
}

type MatchLookTickMsg struct{}

func MatchLookTick() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
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
			select {
			case s, ok := <-m.opC:
				if !ok {
					log.Warn("Tried to read from a closed opC on a MatchLook msg")
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

func (m MultiplayerGame) View() string {
	m.setState()
	switch m.mstate {
	case msLooking:
		return "looking for match"
	case msRunning:
		// TODO: view ours & op's board
		return lipgloss.JoinHorizontal(lipgloss.Left, m.game.View(), RenderBoard(m.opBoard(), defaultBoardStyle()))
	case msCanceled:
		return "match canceled"
	default:
		log.Warn("invalid match state in multiplayergame view function")
		return "something went wrong!"
	}
}
