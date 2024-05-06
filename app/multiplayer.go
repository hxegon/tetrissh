package app

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	opC       <-chan *MultiplayerSession // Should only be read once
	mstate    matchState
}

func NewMultiplayer() *MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int

	session := &MultiplayerSession{
		ctx:   ctx,
		board: board,
	}

	opC := session.requestMatch()

	game := &MultiplayerGame{
		session: session,
		cancel:  cancel,
		opC:     opC,
	}

	return game
}

func (m *MultiplayerGame) close() {
	log.Debug("Closing game")
	m.cancel()
	// drop pointers
	m.opSession.board = nil
	m.session.board = nil
}

// Return enum checking if the game is looking for a match, running or canceled
// sets and returns matchState, but m.mstate shouldn't be accessed other than through here
func (m *MultiplayerGame) state() matchState {
	oldState := m.mstate
	newState := m.mstate

	if m.opSession == nil {
		if oldState == msRunning {
			newState = msCanceled
		}
	} else {
		select {
		// Check if either session in the match is canceled
		case <-m.opSession.done():
			log.Info("opsession canceled, setting mstate to canceled")
			newState = msCanceled
		case <-m.session.done(): // Shouldn't really be reached but just in case
			newState = msCanceled
		default: // If not, make sure mstate is running
			log.Info("default case for state")
			if oldState == msLooking {
				log.Info("Setting multiplayer game to running")
				newState = msRunning
			}
		}
	}

	m.mstate = newState
	return newState
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
	m.state()

	switch msg := msg.(type) {
	case MatchLookTickMsg:
		// Continue refreshing if there's no match found
		log.Infof("MatchLookTickMsg! mstate %v", m.mstate.String())
		if m.mstate == msLooking {
			select {
			case s, ok := <-m.opC:
				if !ok {
					log.Warn("Tried to read from a closed opC on a MatchLook msg")
					// TODO: Return error msg cmd here
					return m, nil
				} else {
					log.Info("Setting opSession")
					m.opSession = s
					return m, nil
				}
			default:
				return m, MatchLookTick()
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			log.Info("Canceling multiplayer ctx")
			m.close()
			return m, DeactivateCmd
		}
	}
	return m, cmd
}

func (m MultiplayerGame) View() string {
	state := m.state()
	switch state {
	case msLooking:
		return "looking for match"
	case msRunning:
		// TODO: view ours & op's board
		return "Found match"
	case msCanceled:
		return "match canceled"
	default:
		log.Warn("invalid match state in multiplayergame view function")
		return "something went wrong!"
	}
}
