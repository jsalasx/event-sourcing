package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"account-cmd/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventStore struct {
	collection *mongo.Collection
}

func NewMongoEventStore(uri, db, collection string) (*MongoEventStore, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	col := client.Database(db).Collection(collection)
	return &MongoEventStore{collection: col}, nil
}

func (s *MongoEventStore) Save(aggregateID string, events []domain.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Len Events:", len(events), "for aggregate ID:", aggregateID)
	docs := make([]interface{}, 0, len(events))
	for _, e := range events {
		raw, _ := json.Marshal(e)

		var eventMap bson.M
		_ = json.Unmarshal(raw, &eventMap) // convierte a un mapa BSON válido

		docs = append(docs, map[string]interface{}{
			"aggregate_id": aggregateID,
			"event_type":   e.EventType(),
			"version":      e.GetVersion(),
			"data":         eventMap,
			"timestamp":    e.Timestamp(),
		})
	}
	_, err := s.collection.InsertMany(ctx, docs)
	return err
}

func (s *MongoEventStore) Load(aggregateID string) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"aggregate_id": aggregateID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []domain.Event

	for cursor.Next(ctx) {
		var doc struct {
			EventType string                 `bson:"event_type"`
			Data      map[string]interface{} `bson:"data"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		// Convertimos el mapa a JSON para usar el registry
		raw, _ := json.Marshal(doc.Data)
		evt, err := domain.BuildEvent(doc.EventType, raw)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *MongoEventStore) LoadWithVersion(aggregateID string, version uint64) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"aggregate_id": aggregateID,
		"version":      bson.M{"$gt": version},
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No hay eventos posteriores a la versión especificada
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []domain.Event

	for cursor.Next(ctx) {
		var doc struct {
			EventType string                 `bson:"event_type"`
			Data      map[string]interface{} `bson:"data"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		raw, _ := json.Marshal(doc.Data)
		evt, err := domain.BuildEvent(doc.EventType, raw)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
