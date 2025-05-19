package frontend

import(
	"net/http"

	"github.com/krishnaGauss/load-balancer/serverpool"
	
)

type LoadBalancer interface{
	Serve(http.ResponseWriter, *http.Request)
}

type loadBalancer struct{
	serverPool serverpool.ServerPool 
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request){
	peer := lb.serverPool.GetNextValidPeer()
	if peer!=nil{
		peer.Serve(w,r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}