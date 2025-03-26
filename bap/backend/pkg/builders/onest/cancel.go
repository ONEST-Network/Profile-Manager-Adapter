package onest


import (
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
)

func BuildBPPCancelJobRequest(payload cancelrequest.SeekerCancelPayload, bppId, bppuri string, worker *workerProfile.WorkerProfile) (*cancelrequest.CancelRequest, error) {
	req := cancelrequest.CancelRequest{
		Context: cancelrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "cancel",
			Version:       "2.0.0",
			TransactionID: worker.TransactionID,
			MessageID:     worker.MessageID,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: cancelrequest.Location{
				City: cancelrequest.City{
					Code: payload.Location.City,
				},
				Country: cancelrequest.Country{
					Code: payload.Location.Country,
				},
			},
		},
		Message: cancelrequest.Message{
			OrderID: worker.ApplicantionID[worker.TransactionID],
		},
	}
	return &req, nil
}