package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Target struct {
	Id      primitive.ObjectID `bson:"_id"     json:"id"`
	OwnerId int64              `bson:"ownerId" json:"ownerId"`
	Url     string             `bson:"url"     json:"url"`
}

type FailedTarget struct {
	Id           primitive.ObjectID `bson:"_id"          json:"id"`
	TargetId     primitive.ObjectID `bson:"targetId"     json:"targetId"`
	StartedAt    *time.Time         `bson:"startedAt"    json:"startedAt"`
	LastFailedAt *time.Time         `bson:"lastFailedAt" json:"lastFailedAt"`
	Resolved     bool               `bson:"resolved"     json:"resolved"`
}

func (ms *mongoStorage) NewTarget(ownerId int64, url string) error {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.TargetsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{primitive.E{Key: "url", Value: url}, primitive.E{Key: "ownerId", Value: ownerId}}
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				primitive.E{Key: "url", Value: url},
				primitive.E{Key: "ownerId", Value: ownerId},
			},
		},
	}
	options := options.FindOneAndUpdate().SetUpsert(true)

	res := collection.FindOneAndUpdate(ctx, filter, update, options)
	return res.Err()
}

func (ms *mongoStorage) setTarget(targetId primitive.ObjectID, status bool) error {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.StatusCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	ts := time.Now()
	filter := bson.D{primitive.E{Key: "targetId", Value: targetId}, primitive.E{Key: "resolved", Value: status}}
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				primitive.E{Key: "resolved", Value: status},
				primitive.E{Key: "lastFailedAt", Value: ts},
			},
		},
		primitive.E{
			Key: "$setOnInsert", Value: bson.D{
				primitive.E{Key: "startedAt", Value: ts},
			},
		},
	}
	options := options.FindOneAndUpdate().SetUpsert(true)

	res := collection.FindOneAndUpdate(ctx, filter, update, options)
	return res.Err()
}

func (ms *mongoStorage) IsResolved(targetId primitive.ObjectID) (bool, error) {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.StatusCollection)
	ctx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{primitive.E{Key: "targetId", Value: targetId}}
	options := options.FindOne().SetSort(
		bson.D{
			primitive.E{Key: "_id", Value: -1},
		},
	)

	res := collection.FindOne(ctx, filter, options)

	var targetStatus FailedTarget

	if err := res.Decode(&targetStatus); err != nil {
		// If document doesn't exist then it means that we haven't had any failures yet
		if err == mongo.ErrNoDocuments {
			return true, nil
		}

		return true, err
	}

	return targetStatus.Resolved, nil
}

func (ms *mongoStorage) FailTarget(targetId primitive.ObjectID) error {
	return ms.setTarget(targetId, false)
}

func (ms *mongoStorage) ResolveTarget(targetId primitive.ObjectID) error {
	return ms.setTarget(targetId, true)
}

func (ms *mongoStorage) fetchByCriteria(criteria ...primitive.E) ([]Target, error) {
	db := ms.client.Database(ms.config.Storage.Database)
	collection := db.Collection(ms.config.Storage.TargetsCollection)

	countCtx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	filter := bson.D{}
	for _, c := range criteria {
		filter = append(filter, c)
	}

	count, err := collection.CountDocuments(countCtx, filter)

	if err != nil {
		return nil, err
	}

	selectCtx, cancel := context.WithTimeout(context.Background(), ms.config.Storage.WriteTimeoutDuration)
	defer cancel()

	cursor, err := collection.Find(selectCtx, filter)

	if err != nil {
		return nil, err
	}

	result := make([]Target, 0, count)

	for cursor.Next(context.Background()) {
		var target Target

		if err := cursor.Decode(&target); err != nil {
			ms.logger.Errorf("Couldn't decode document into target. Err - %s", err.Error())
		}

		result = append(result, target)
	}

	return result, nil
}

func (ms *mongoStorage) FetchTargetsByOwner(userId int64) ([]Target, error) {
	return ms.fetchByCriteria(primitive.E{Key: "ownerId", Value: userId})
}

func (ms *mongoStorage) FetchTargets() ([]Target, error) {
	return ms.fetchByCriteria()
}
