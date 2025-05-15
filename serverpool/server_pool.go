package serverpool

import(
	"fmt"
	"sync"
	"github.com/krishnaGauss/load-balancer/backend"
)


type ServerPool interface{
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	addBackend(backend.Backend)
	GetServerPoolSize() int
}

type roundRobinServerPool struct{
	backends []backend.Backend
	mux *sync.RWMutex
	current int
}

func (s *roundRobinServerPool) Rotate(backend.Backend){
	s.mux.Lock()
	s.current = (s.current+1)
	
}