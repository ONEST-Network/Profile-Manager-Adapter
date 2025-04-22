package workerProfile

import (
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	workerProfileDb "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	workerProfilePayload "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/workerProfile"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Interface interface {
	AddWorkerProfile(business *workerProfilePayload.AddWorkerProfileRequest) error
}

type WorkerProfile struct {
	clients *clients.Clients
}

func NewWorkerProfile(clients *clients.Clients) Interface {
	return &WorkerProfile{
		clients: clients,
	}
}

func (b *WorkerProfile) AddWorkerProfile(payload *workerProfilePayload.AddWorkerProfileRequest) error {
	logrus.Infof("[Request]: Received request to add a new worker profile: %s", payload.Name)

	workers, err := b.clients.WorkerProfileClient.ListWorkerProfile(bson.D{{Key: "id", Value: payload.ID}})
	if err != nil {
		logrus.Errorf("Failed to get worker profile with id %s, %v", payload.ID, err)
		return fmt.Errorf("failed to get worker profile with id %s, %v", payload.ID, err)
	}

	if len(workers) > 0 {
		logrus.Errorf("worker profile with id %s already exists", payload.ID)
		return fmt.Errorf("worker profile with id %s already exists", payload.ID)
	}

	var worker = &workerProfileDb.WorkerProfile{
		ID:             payload.ID,
		Name:           payload.Name,
		Phone:          payload.Phone,
		Email:          payload.Email,
		Location:       payload.Location,
		Age:            payload.Age,
		Gender: 	   payload.Gender,
		PreferredLanguages: payload.PreferredLanguages,
		PreferredJobRoles:  payload.PreferredJobRoles,
		Experience:    payload.Experience,
		Skills:        payload.Skills,
		Certifications: payload.Certifications,
		Credentials:    payload.Credentials,
		ActiveJobApplications: make(map[string]workerProfileDb.ActiveJobApplications),
	}

	if err := b.clients.WorkerProfileClient.CreateWorkerProfile(worker); err != nil {
		logrus.Errorf("Failed to create %s worker profile, %v", payload.ID, err)
		return fmt.Errorf("failed to create %s worker profile, %v", payload.ID, err)
	}

	return nil
}
