package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/krishnaGauss/load-balancer/backend"
	"github.com/krishnaGauss/load-balancer/frontend"
	"github.com/krishnaGauss/load-balancer/serverpool"
	// "github.com/krishnaGauss/load-balancer/utils"
)

func main() {
	logger := utils.InitLogger()
	defer logger.Sync()

	config, err := utils.GetLBConfig()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPool, err := serverpool.NewServerPool(utils.GetLBStrategy(config.Strategy))
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	loadBalancer := frontend.NewLoadBalancer(serverPool)

	for _, u := range config.Backends {
		endpoint, err := url.Parse(u)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("URL", u))
		}

		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := backend.NewBackend(endpoint, rp)
		rp.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			logger.Error("error handling the request",
				zap.String("host", endpoint.Host),
				zap.Error(e),
			)
			backendServer.SetAlive(false)

			if !frontend.AllowRetry(request) {
				utils.Logger.Info(
					"Max retry attempts reached, terminating",
					zap.String("address", request.RemoteAddr),
					zap.String("path", request.URL.Path),
				)
				http.Error(writer, "Service not available", http.StatusServiceUnavailable)
				return
			}

			logger.Info(
				"Attempting retry",
				zap.String("address", request.RemoteAddr),
				zap.String("URL", request.URL.Path),
				zap.Bool("retry", true),
			)
			loadBalancer.Serve(
				writer,
				request.WithContext(
					context.WithValue(request.Context(), frontend.RETRY_ATTEMPTED, true),
				),
			)
		}

		serverPool.AddBackend(backendServer)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(loadBalancer.Serve),
	}

	go serverpool.LauchHealthCheck(ctx, serverPool)

	go func() {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Info(
		"Load Balancer started",
		zap.Int("port", config.Port),
	)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("ListenAndServe() error", zap.Error(err))
	}
}