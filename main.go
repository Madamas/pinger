package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"pinger/packages/config"
	"pinger/packages/logger"
	"pinger/packages/pinger"
	"pinger/packages/server"
	"pinger/packages/storage"
	"pinger/packages/telegram"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			server.AsRoute(server.NewAliveHandler),
			fx.Annotate(
				server.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			server.NewHTTPServer,
			config.Populate,
			logger.NewZap,
			pinger.NewPinger,
			pinger.NewClient,
			storage.NewClient,
			storage.NewStorage,
			telegram.NewBot,
			telegram.NewNotifier,
			telegram.NewUpdateChannel,
			telegram.NewListener,
		),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(pinger.Pinger) {}),
		fx.Invoke(func(telegram.Listener) {}),
		fx.StartTimeout(1*time.Second),
		fx.StopTimeout(1*time.Second),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	log.Println("App starting...")
	err := app.Start(startCtx)

	if err != nil {
		panic(err)
	}

	log.Println("Application successfully started")
	sigs := app.Done()
	sig := <-sigs

	log.Println("\n Received signal: ", sig)
	log.Println("Exiting in " + app.StopTimeout().String())
	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()

	fmt.Println("Stopping application...")

	log.Fatal(app.Stop(stopCtx))
}
