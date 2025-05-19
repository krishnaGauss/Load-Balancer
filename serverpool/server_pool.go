package serverpool

import (
	"context"
	"sync"

	"github.com/krishnaGauss/load-balancer/backend"
)

type ServerPool interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
	GetServerPoolSize() int
}

type roundRobinServerPool struct {
	backends []backend.Backend
	mux      *sync.RWMutex
	current  int
}

func (s *roundRobinServerPool) Rotate() backend.Backend {
	s.mux.Lock()
	s.current = (s.current + 1) % s.GetServerPoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func (s *roundRobinServerPool) GetNextValidPeer() backend.Backend {
	for i := 0; i < s.GetServerPoolSize(); i++ {
		nextPeer := s.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil

}

func (s *roundRobinServerPool) GetBackends() []backend.Backend{
	return s.backends
}

func (s *roundRobinServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func (s *roundRobinServerPool) AddBackend(b backend.Backend) {
	s.backends = append(s.backends, b)
}


func HealthCheck(ctx context.Context, s ServerPool){
	aliveChannel := make(chan bool, 1)
	for _,b:=range s.GetBackends(){
		b:=b
	}

	//defining context with timeout
}
