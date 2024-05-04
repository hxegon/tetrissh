package app

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var matchReqC = make(chan matchReq)

type MultiplayerSession struct {
	ctx   context.Context
	board *[][]int
	mx    *sync.RWMutex
}

type matchReq struct {
	session   *MultiplayerSession
	opSession *MultiplayerSession
}

// Return the done channel from the request's context
func (m matchReq) done() <-chan struct{} {
	return m.session.ctx.Done()
}

// On a loop, match requests. Meant to be used in a goroutine in main
// TODO: pass in a context
func MatchSessions() {
	var lastReq *matchReq

	for {
		nextReq := <-matchReqC

		if lastReq == nil {
			lastReq = &nextReq
		} else {
			select {
			case <-lastReq.done():
				log.Info("Multiplayer request canceled")
				*lastReq = nextReq
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

func NewMultiplayer() MultiplayerGame {
	ctx, cancel := context.WithCancel(context.Background())
	var board *[][]int
	var mx *sync.RWMutex

	session := MultiplayerSession{
		ctx:   ctx,
		board: board,
		mx:    mx,
	}

	var opSession *MultiplayerSession

	matchReqC <- matchReq{
		session:   &session,
		opSession: opSession,
	}

	return MultiplayerGame{
		session:   &session,
		opSession: opSession,
		cancel:    cancel,
	}
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
