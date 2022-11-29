package pinger

import (
	"context"
	"fmt"
	"net/http"
	"pinger/packages/config"
	"pinger/packages/storage"
	"pinger/packages/telegram"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Pinger struct {
	config   config.Pinger
	quitter  chan int
	client   *http.Client
	targets  []storage.Target
	logger   *zap.SugaredLogger
	storage  storage.Storage
	notifier telegram.Notifier
}

func NewPinger(
	client *http.Client,
	conf config.Config,
	logger *zap.SugaredLogger,
	st storage.Storage,
	lc fx.Lifecycle,
	notifier telegram.Notifier,
) (Pinger, error) {
	targets, err := st.FetchTargets()

	if err != nil {
		return Pinger{}, err
	}

	p := Pinger{
		config:   conf.Pinger,
		quitter:  make(chan int, 1),
		client:   client,
		targets:  targets,
		logger:   logger,
		storage:  st,
		notifier: notifier,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if conf.Pinger.Enabled {
				p.run()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.quitter <- 1
			return nil
		},
	})

	return p, nil
}

func (p Pinger) run() {
	p.logger.Infoln("Starting pinger")

	t := time.NewTicker(p.config.IntervalDuration)
	tt := time.NewTicker(p.config.ReloadIntervalDuration)

	go func() {
		for {
			select {
			case <-tt.C:
				p.reloadTargets()
			case <-t.C:
				p.rotateRequest()
			case <-p.quitter:
				p.logger.Info("Stopping pinger")
				return
			}
		}
	}()
}

func (p *Pinger) reloadTargets() {
	newTargets, err := p.storage.FetchTargets()
	if err != nil {
		p.logger.Errorf("Failed to reload targets for pinger, err - %s", err.Error())
	}
	p.targets = newTargets
}

func (p Pinger) rotateRequest() {
	for _, target := range p.targets {
		errorStatus := false
		
		p.logger.Infof("Sending request to %s", target.Url)

		resp, err := p.client.Get(string(target.Url))

		if err != nil {
			p.logger.Errorf("Receiver error from target, %s, Error: %s", target.Url, err.Error())
			errorStatus = true
		} else {
			p.logger.Infof("Received status code %d", resp.StatusCode)
			if resp.StatusCode > 300 {
				p.logger.Info("Status code is higher than 300. Perceiving target as errored")
				errorStatus = true
			}
		}

		isOk, err := p.storage.IsResolved(target.Id)

		if err != nil {
			p.logger.Errorf("Couldn't check if target %s was ok or not. Err - %s", target.Id.Hex(), err.Error())
			errorStatus = true
		}

		if errorStatus {
			if isOk {
				if err := p.storage.FailTarget(target.Id); err != nil {
					p.logger.Errorf("Couldn't mark target %s as failed, err - %s", target.Id.Hex(), err.Error())
				}

				p.notifier.Notify(fmt.Sprintf("Target %s has failed", target.Url), target.OwnerId)
			}
		} else {
			if !isOk {
				if err := p.storage.ResolveTarget(target.Id); err != nil {
					p.logger.Errorf("Couldn't resolve target %s as failed, err - %s", target.Id.Hex(), err.Error())
				}

				p.notifier.Notify(fmt.Sprintf("Target %s was resolved", target.Url), target.OwnerId)
			}
		}
	}
}
