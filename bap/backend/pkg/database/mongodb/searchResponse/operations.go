package searchresponse

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetSearchJobResponse(id string) (*SearchJobResponse, error)
	CreateSearchJobResponse(jobResponse *SearchJobResponse) error
	ListSearchJobResponse(query bson.D) ([]SearchJobResponse, error)
	DeleteSearchJobResponse(jobResponse, name string) error
	UpdateSearchJobResponse(query, update bson.D) error
}

type Dao struct {
	Collection *mongo.Collection
}

func NewSearchJobResponseDao(collection *mongo.Collection) *Dao {
	return &Dao{
		Collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetSearchJobResponse(id string) (*SearchJobResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var jobResponse SearchJobResponse
	if err := d.Collection.FindOne(ctx, bson.D{{Key: "transaction_id", Value: id}}).Decode(&jobResponse); err != nil {
		return nil, err
	}

	return &jobResponse, nil
}

func (d *Dao) CreateSearchJobResponse(jobResponse *SearchJobResponse) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()
    _, err := d.Collection.InsertOne(ctx, jobResponse)
    return err
}

func (d *Dao) ListSearchJobResponse(query bson.D) ([]SearchJobResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    cursor, err := d.Collection.Find(ctx, query)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var jobResponses []SearchJobResponse
    if err = cursor.All(ctx, &jobResponses); err != nil {
        return nil, err
    }

    return jobResponses, nil
}

func (d *Dao) DeleteSearchJobResponse(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.Collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}

func (d *Dao) UpdateSearchJobResponse(query, update bson.D) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.Collection.UpdateOne(ctx, query, update)
    return err
}
