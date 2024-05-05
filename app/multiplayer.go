package app

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var matchReqC = make(chan *MultiplayerGame)

type MultiplayerSession struct {
	ctx   context.Context
	board *[][]int
	mx    *sync.RWMutex
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
				log.Info("Multiplayer request canceled")
				lastReq = nextReq
			default:
				log.Info("Exchanging multiplayer session pointers")
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

type MultiplayerGame struct {
	cancel    context.CancelFunc
	session   *MultiplayerSession
	opSession *MultiplayerSession
	mstate    matchState
}

func NewMultiplayer() *MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int
	var mx *sync.RWMutex

	session := MultiplayerSession{
		ctx:   ctx,
		board: board,
		mx:    mx,
	}

	game := &MultiplayerGame{
		// opSession is left as an uninitialized pointer to be used later by the matchmaking goroutine
		session: &session,
		cancel:  cancel,
	}

	matchReqC <- game

	return game
}

// sets and returns matchState, but m.mstate shouldn't be accessed other than through here
func (m *MultiplayerGame) setState() matchState {
	if m.opSession == nil {
		if m.mstate == msRunning {
			m.mstate = msCanceled
		}
	} else {
		select {
		case <-m.opSession.done():
			m.opSession = nil
			m.cancel()
			m.mstate = msCanceled
		case <-m.session.done(): // Shouldn't really be reached but just in case
			m.opSession = nil
			m.mstate = msCanceled
		default:
			if m.mstate == msLooking {
				m.mstate = msRunning
			}
		}
	}

	return m.mstate
}

func (m MultiplayerGame) Init() tea.Cmd {
	return nil
}

// TODO: Start ticker to refresh view when looking for match
func (m MultiplayerGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// state := m.setState()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			log.Info("Canceling multiplayer ctx")
			m.cancel()
			return m, DeactivateCmd
		}
	}
	return m, nil
}

func (m MultiplayerGame) View() string {
	switch m.setState() {
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
