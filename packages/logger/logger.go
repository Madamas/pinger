package logger

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewZap(
	lc fx.Lifecycle,
) *zap.SugaredLogger {
	logger, _ := zap.NewProduction()

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return logger.Sync()
		},
	})

	return logger.Sugar()
}
