package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"pinger/packages/config"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPServer(
	config config.Config,
	lc fx.Lifecycle,
	mux *http.ServeMux,
	logger *zap.SugaredLogger,
) *http.Server {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting HTTP server at %s\n", server.Addr)
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
