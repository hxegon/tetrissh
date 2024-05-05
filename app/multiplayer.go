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

// On a loop, match requests. Meant to be used in a goroutine in main
func MatchMultiplayerGames() {
	var lastReq *MultiplayerGame

	for {
		nextReq := <-matchReqC

		if lastReq == nil {
			lastReq = nextReq
		} else {
			select {
			case <-lastReq.done(): // Has the last request been canceled?
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

func (m MultiplayerGame) done() <-chan struct{} {
	return m.session.ctx.Done()
}

func (m MultiplayerGame) Init() tea.Cmd {
	return nil
}

func (m MultiplayerGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
