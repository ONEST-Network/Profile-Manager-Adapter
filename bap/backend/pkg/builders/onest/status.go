package onest

import (
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
)

func BuildBPPStatusJobRequest(payload statusrequest.SeekerStatusPayload, bppId, bppuri string, worker *workerProfile.WorkerProfile, jobCityCode, jobCountryCode string) (*statusrequest.StatusRequest, error) {
	req := statusrequest.StatusRequest{
		Context: statusrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "status",
			Version:       "2.0.0",
			TransactionID: worker.TransactionID,
			MessageID:     worker.MessageID,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: statusrequest.Location{
				City: statusrequest.City{
					Code: jobCityCode,
				},
				Country: statusrequest.Country{
					Code: jobCountryCode,
				},
			},
		},
		Message: statusrequest.Message{
			Order: statusrequest.Order{
				ID: payload.ApplicationID,
			},
		},
	}
	return &req, nil
}