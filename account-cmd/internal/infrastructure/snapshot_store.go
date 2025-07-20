package infrastructure

import (
	"account-cmd/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSnapshotStore struct {
	collection *mongo.Collection
}

func NewMongoSnapshotStore(uri, db, collection string) (*MongoSnapshotStore, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	col := client.Database(db).Collection(collection)
	return &MongoSnapshotStore{collection: col}, nil
}

func (s *MongoSnapshotStore) UpsertSnapshot(snapshot domain.BankAccountSnapshot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	filter := bson.M{"aggregate_id": snapshot.AggregateID}
	update := bson.M{
		"$set": bson.M{
			"aggregate_id": snapshot.AggregateID,
			"version":      snapshot.Version,
			"state":        snapshot.State,
			"created_at":   snapshot.CreatedAt,
		},
	}

	_, err := s.collection.UpdateOne(
		ctx,
		filter,
		update,
		options.Update().SetUpsert(true), // inserta si no existe
	)

	return err
}

func (s *MongoSnapshotStore) LoadSnapshot(aggregateID string) (*domain.BankAccountSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var snapshot domain.BankAccountSnapshot
	err := s.collection.FindOne(ctx, bson.M{"aggregate_id": aggregateID}).Decode(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No snapshot found
		}
		return nil, err // Other error
	}

	return &snapshot, nil
}
