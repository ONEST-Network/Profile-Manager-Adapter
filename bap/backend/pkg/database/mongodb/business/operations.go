package business

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetBusiness(id string) (*Business, error)
	CreateBusiness(business *Business) error
	ListBusiness(query bson.D) ([]Business, error)
	DeleteBusiness(businessID, name string) error
	UpdateBusiness(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewBusinessDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetBusiness(id string) (*Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var business Business
	if err := d.collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&business); err != nil {
		return nil, err
	}

	return &business, nil
}
