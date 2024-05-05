package app

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var matchReqC = make(chan *MultiplayerGame)

type MultiplayerSession struct {
	ctx   context.Context
	board *[][]int
}

func (m MultiplayerSession) done() <-chan struct{} {
	return m.ctx.Done()
}

// On a loop, match requests. Meant to be used in a goroutine in main
func MatchMultiplayerGames() {
	var lastReq *MultiplayerGame

	for {
		nextReq := <-matchReqC

		if lastReq == nil {
			lastReq = nextReq
		} else {
			select {
			case <-lastReq.session.done(): // Has the last request been canceled?
				log.Info("match request canceled")
				lastReq = nextReq
			default:
				log.Info("Exchanging match pointers")
				lastReq.opSession = nextReq.session
				nextReq.opSession = lastReq.session

				lastReq = nil
			}
		}
	}
}

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
	mstate    matchState
}

func NewMultiplayer() *MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int

	session := &MultiplayerSession{
		ctx:   ctx,
		board: board,
	}
	var opSession *MultiplayerSession

	game := &MultiplayerGame{
		session:   session,
		opSession: opSession,
		cancel:    cancel,
	}

	matchReqC <- game

	return game
}

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
			m.opSession = nil
			m.cancel()
			newState = msCanceled
		case <-m.session.done(): // Shouldn't really be reached but just in case
			m.opSession = nil
			newState = msCanceled
		default: // If not, make sure mstate is running
			if oldState == msLooking {
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
		log.Info("MatchLookTickMsg")
		if m.mstate == msLooking {
			cmd = MatchLookTick()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			log.Info("Canceling multiplayer ctx")
			m.cancel()
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
		return "Found match"
	case msCanceled:
		return "match canceled"
	default:
		log.Warn("invalid match state in multiplayergame view function")
		return "something went wrong!"
	}
}
