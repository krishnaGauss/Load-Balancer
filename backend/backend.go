package backend

import (
	"sync"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend interface{
	SetAlive(bool)
	IsAlive() bool
	GetURL() *url.URL
	GetActiveConnections() int
	Serve(http.ResponseWriter, *http.Request)
}

type backend struct{
	url *url.URL
	alive bool
	mux sync.RWMutex
	connections int
	reverseProxy *httputil.ReverseProxy
}