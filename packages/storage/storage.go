package storage

import (
	"pinger/packages/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Storage interface {
	FetchTargets() ([]Target, error)
	FetchTargetsByOwner(ownerId int64) ([]Target, error)
	NewTarget(ownerId int64, url string) error
	IsResolved(targetId primitive.ObjectID) (bool, error)
	FailTarget(targetId primitive.ObjectID) error
	ResolveTarget(targetId primitive.ObjectID) error

	AddUser(id int64) (User, error)
	GetStatus(userId int64) (Status, error)
	SetStatus(userId int64, status Status) error
}

type mongoStorage struct {
	client *mongo.Client
	config config.Config
	logger *zap.SugaredLogger
}

func NewStorage(
	mongo *mongo.Client,
	config config.Config,
) Storage {
	return &mongoStorage{
		client: mongo,
		config: config,
	}
}
