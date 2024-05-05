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

type MultiplayerGame struct {
	cancel    context.CancelFunc
	session   *MultiplayerSession
	opSession *MultiplayerSession
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

func (m MultiplayerGame) Init() tea.Cmd {
	return nil
}

func (m MultiplayerGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Drop opSession pointer if it exists but was canceled
	// TODO: Notify user that the opSession was canceled
	if m.opSession != nil {
		select {
		case <-m.opSession.done():
			m.opSession = nil
		default:
		}
	}

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
	if m.opSession == nil {
		return "looking for match"
	} else {
		return "found match"
	}
}
