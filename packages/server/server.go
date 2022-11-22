package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"pinger/packages/config"
	"time"

	"go.uber.org/fx"
)

func NewHTTPServer(
	config config.Config,
	lc fx.Lifecycle,
	mux *http.ServeMux,
) *http.Server {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting HTTP server at", server.Addr)
			ln, err := net.Listen("tcp", server.Addr)
			fmt.Println("Hmmm")
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
