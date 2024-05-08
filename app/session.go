package app

import (
	"context"
	"sync"

	"github.com/charmbracelet/log"
)

type MultiplayerSession struct {
	ctx   context.Context
	board *[][]int
	mx    *sync.RWMutex
}

func (m *MultiplayerSession) done() <-chan struct{} {
	return m.ctx.Done()
}

// Returns the current board state and an ok bool
// ok bool will fail if the board is inaccessible. (Session is closed or the board pointer is nil)
// Thread safe, blocks for a mutex lock
func (m *MultiplayerSession) Board() (board [][]int, ok bool) {
	select {
	case <-m.done():
		ok = false
	default:
		if m.board == nil {
			log.Warn("session context wasn't closed but board pointer was nil")
			ok = false
		} else {
			m.mx.RLock()
			defer m.mx.RUnlock()

			board = *m.board
			ok = true
		}
	}
	return board, ok
}

// Thread safe setter. Blocks for mutex.
func (m *MultiplayerSession) SetBoard(b [][]int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.board = &b
}

type matchReq struct {
	session *MultiplayerSession
	opC     chan<- *MultiplayerSession
}

var matchReqC = make(chan matchReq)

// Request a match and return a recieving channel that the opponent's session will be returned through
func (s *MultiplayerSession) requestMatch() <-chan *MultiplayerSession {
	opC := make(chan *MultiplayerSession, 1) // Don't want to block matchmaking when sending

	matchReqC <- matchReq{
		session: s,
		opC:     opC,
	}

	return opC
}

// On a loop, match requests. Meant to be used in a goroutine in main
func MatchMultiplayerGames() {
	var lastReq *matchReq

	for {
		nextReq := <-matchReqC

		if lastReq == nil {
			lastReq = &nextReq
		} else {
			select {
			case <-lastReq.session.done(): // Has the last request been canceled?
				log.Debug("match request canceled")
				close(lastReq.opC)
				*lastReq = nextReq
			case <-nextReq.session.done():
				// Skip this request if context is canceled
				continue
			default:
				log.Debug("Exchanging match requests")
				// FIXME: The sending channels shouldn't be filled anywhere else, but we should still check/handle it if they are
				// otherwise this will hang matchmaking. Also, handle panics here.
				lastReq.opC <- nextReq.session
				close(lastReq.opC)
				nextReq.opC <- lastReq.session
				close(nextReq.opC)

				lastReq = nil
			}
		}
	}
}
