package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"shared/ports"
)

type MongoAccountRepository struct {
	collection *mongo.Collection
}

func NewMongoAccountRepository(uri, db, collection string) (*MongoAccountRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	col := client.Database(db).Collection(collection)
	return &MongoAccountRepository{collection: col}, nil
}

func (r *MongoAccountRepository) Upsert(account ports.AccountRead) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"account_id": account.AccountID},
		bson.M{"$set": account},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *MongoAccountRepository) UpdateBalance(accountID string, delta float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"account_id": accountID},
		bson.M{"$inc": bson.M{"balance": delta}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *MongoAccountRepository) FindByID(accountID string) (*ports.AccountRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result ports.AccountRead
	err := r.collection.FindOne(ctx, bson.M{"account_id": accountID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
