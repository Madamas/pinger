package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Status string

const (
	STATUS_UNDEFINED Status = "undefined"

	STATUS_INITIAL    Status = "initial"
	STATUS_NEW_TARGET Status = "new_target"
)

func (s Status) IsError() bool {
	return s == STATUS_UNDEFINED
}

type User struct {
	Id        int64      `bson:"_id"     json:"id"`
	CreatedAt *time.Time `bson:"createdAt" json:"createdAt"`
	Status    Status     `bson:"status" json:"status"`
}

func (ms *mongoStorage) AddUser(id int64) (User, error) {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{
			Key: "$setOnInsert", Value: bson.D{
				primitive.E{Key: "status", Value: STATUS_INITIAL},
				primitive.E{Key: "createdAt", Value: time.Now()},
			},
		},
	}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	res := collection.FindOneAndUpdate(ctx, filter, update, options)
	user := User{}

	if res.Err() != nil {
		return user, res.Err()
	}

	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (ms *mongoStorage) GetStatus(userId int64) (Status, error) {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}

	options := options.FindOne()
	res := collection.FindOne(ctx, filter, options)
	user := User{}

	if res.Err() != nil {
		return STATUS_UNDEFINED, res.Err()
	}

	if err := res.Decode(&user); err != nil {
		return STATUS_UNDEFINED, err
	}

	return user.Status, nil
}

func (ms *mongoStorage) SetStatus(userId int64, status Status) error {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				primitive.E{Key: "status", Value: status},
			},
		},
		primitive.E{
			Key: "$setOnInsert", Value: bson.D{
				primitive.E{Key: "createdAt", Value: time.Now()},
			},
		},
	}
	options := options.FindOneAndUpdate().SetUpsert(true)

	res := collection.FindOneAndUpdate(ctx, filter, update, options)

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
