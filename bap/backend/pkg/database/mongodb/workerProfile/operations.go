package workerProfile

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DaoInterface interface {
	GetWorkerProfile(id string) (*WorkerProfile, error)
	CreateWorkerProfile(worker *WorkerProfile) error
	ListWorkerProfile(query bson.D) ([]WorkerProfile, error)
	DeleteWorkerProfile(worker, name string) error
	UpdateWorkerProfile(query, update bson.D) error
}

type Dao struct {
	collection *mongo.Collection
}

func NewWorkerDao(collection *mongo.Collection) *Dao {
	return &Dao{
		collection: collection,
	}
}

const dbTimeout = 10 * time.Second

func (d *Dao) GetWorkerProfile(id string) (*WorkerProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var worker WorkerProfile
	if err := d.collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

func (d *Dao) CreateWorkerProfile(worker *WorkerProfile) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()
    // adds the credential short and long description from credentials name
    if err := parseWorkProfileCredentials(worker); err != nil {
        return err
    }
    _, err := d.collection.InsertOne(ctx, worker)
    return err
}

func (d *Dao) ListWorkerProfile(query bson.D) ([]WorkerProfile, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    cursor, err := d.collection.Find(ctx, query)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var workers []WorkerProfile
    if err = cursor.All(ctx, &workers); err != nil {
        return nil, err
    }

    return workers, nil
}

func (d *Dao) DeleteWorkerProfile(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}

func (d *Dao) UpdateWorkerProfile(query, update bson.D) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    _, err := d.collection.UpdateOne(ctx, query, update)
    return err
}

func parseWorkProfileCredentials(worker *WorkerProfile) error {
    for i, cred := range worker.Credentials {
        credential, err := GetWorkProfileCredentialShortandLongDesc(cred.Name)
        if err != nil {
            return err
        }
        worker.Credentials[i].ShortDesc = credential["short_desc"]
        worker.Credentials[i].LongDesc = credential["long_desc"]
    }
    return nil
}

func GetWorkProfileCredentialShortandLongDesc(credName string) (map[string]string, error) {
	credentialDetails := map[string]map[string]string{
		"AADHAAR_CARD":  {"short_desc": "Aadhaar Card information", "long_desc": "Unique Identification Number"},
		"DRIVING_LICENSE":  {"short_desc": "Driving License information", "long_desc": "License to Drive"},
		"RESUME":  {"short_desc": "Summary of qualifications, work experience, and education.", "long_desc": "A comprehensive document showcasing an individual's career achievements, skills, work history, education, certifications, and professional experience."},
		"PASSPORT":  {"short_desc": "Passport information", "long_desc": "International Travel Document"},
		"VOTER_ID":  {"short_desc": "Voter ID information", "long_desc": "Election Commission ID"},
		"PAN_CARD":  {"short_desc": "PAN Card information", "long_desc": "Permanent Account Number"},
	}

	// Normalize input credential name to handle case sensitivity and trim spaces
	if credentialDetail, exists := credentialDetails[strings.TrimSpace(credName)]; exists {
		return credentialDetail, nil
	}

	return nil, fmt.Errorf("credential details not found for %s", credName)
}
