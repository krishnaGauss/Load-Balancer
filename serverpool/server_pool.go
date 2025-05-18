package serverpool

import (
	"fmt"
	"github.com/krishnaGauss/load-balancer/backend"
	"sync"
)

type ServerPool interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	addBackend(backend.Backend)
	GetServerPoolSize() int
}

type roundRobinServerPool struct {
	backends []backend.Backend
	mux      *sync.RWMutex
	current  int
}

func (s *roundRobinServerPool) Rotate() backend.Backend {
	s.mux.Lock()
	s.current = (s.current + 1) % GetServerPoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func (s *roundRobinServerPool) GetNextValidPeer() backend.Backend {
	for i := 0; i < GetServerPoolSize(); i++ {
		nextPeer := s.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil

}

func (s *roundRobinServerPool)GetServerPoolSize() int{
	return len(s.backends)
}


