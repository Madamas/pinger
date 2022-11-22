package pinger

import (
	"context"
	"fmt"
	"net/http"
	"pinger/packages/config"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Pinger struct {
	interval time.Duration
	quitter  chan int
	client   *http.Client
	targets  []config.PingerTarget
	logger   *zap.SugaredLogger
}

func NewPinger(
	client *http.Client,
	config config.Config,
	logger *zap.SugaredLogger,
	lc fx.Lifecycle,
) Pinger {
	p := Pinger{
		interval: config.Pinger.IntervalDuration,
		quitter:  make(chan int),
		client:   client,
		targets:  config.Pinger.Targets,
		logger:   logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.quitter <- 1
			return nil
		},
	})

	return p
}

func (p Pinger) run() {
	p.logger.Infoln("Starting pinger")

	t := time.NewTicker(p.interval)

	go func() {
		select {
		case <-t.C:
			p.logger.Infof("Pinger tick")
			p.rotateRequest()
		case <-p.quitter:
			p.logger.Info("Stopping pinger")
			return
		}
	}()
}

func (p Pinger) rotateRequest() {
	for _, v := range p.targets {
		url := fmt.Sprintf("http://%s:%d%s", v.Host, v.Port, v.Route)

		p.logger.Infof("Sending request to %s", url)

		resp, err := p.client.Get(url)

		p.logger.Infof("Received status code %d", resp.StatusCode)

		if err != nil {
			p.logger.Errorf("Receiver error from target, %s, Error: %s", v.Host, err.Error())
		}
	}
}
