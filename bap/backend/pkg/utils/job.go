package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/service"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/sirupsen/logrus"

	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
)

type JobSync struct {
	onestService *service.OnestBPPService
	interval    time.Duration
	stopChan    chan struct{}
}

func GetEmptySearchRequest() *searchrequest.SearchRequest {
	return &searchrequest.SearchRequest{
		Context: searchrequest.Context{
			Domain:        "jobs",
			Action:        "search",
			Version:       "1.0.0",
			BapID:         "example-bap",
			BapURI:        "https://example.com/callback",
			TransactionID: fmt.Sprintf("tx-%s", time.Now().Format("20060102150405")),
			MessageID:     fmt.Sprintf("msg-%s", time.Now().Format("20060102150405")),
			Timestamp:     time.Now().Format(time.RFC3339),
			Location: searchrequest.Location{
				City: searchrequest.City{
					Code: "std:080",
				},
				Country: searchrequest.Country{
					Code: "IND",
				},
			},
		},
		Message: searchrequest.Message{
			Intent: searchrequest.Intent{},
		},
	}
}

func NewJobSync(clients *clients.Clients, interval time.Duration) *JobSync {
	return &JobSync{
		onestService: service.NewOnestBPPService(clients),
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

func (j *JobSync) Start() {
	ticker := time.NewTicker(j.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				j.syncJobs()
			case <-j.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
	logrus.Info("Job sync started with interval: ", j.interval)
}

func (j *JobSync) Stop() {
	close(j.stopChan)
	logrus.Info("Job sync stopped")
}

func (j *JobSync) syncJobs() {
	logrus.Info("Starting job sync")
	ctx := context.Background()
	emptyRequest := GetEmptySearchRequest()
	response, err := j.onestService.Search(ctx, emptyRequest)
    if err != nil {
        logrus.Errorf("Failed to sync jobs: %v", err)
        return
    }
    logrus.Infof("Successfully synced jobs, found %d providers", len(response.Message.Catalog.Providers))
}
