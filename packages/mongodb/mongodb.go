package mongodb

import (
	"context"
	"pinger/packages/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewClient(
	config config.Config,
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logger.Infof("Connection url %s", config.Storage.Url)
	clientOptions := options.Client().ApplyURI(config.Storage.Url)
	client, err := mongo.Connect(ctx, clientOptions)

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			return client.Ping(c, nil)
		},
		OnStop: func(c context.Context) error {
			return client.Disconnect(c)
		},
	})

	return client, err
}
