package onest

import (
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
)

func BuildBPPSelectJobRequest(payload selectrequest.SeekerSelectPayload, transactionId, messageId, bppId, bppuri string) (*selectrequest.SelectRequest, error) {
	req := selectrequest.SelectRequest{
		Context: selectrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "select",
			Version:       "2.0.0",
			TransactionID: transactionId,
			MessageID:     messageId,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC(),
			TTL:           "PT30S",
			Location: selectrequest.Location{
				City: selectrequest.City{
					Code: payload.Location.City,
				},
				Country: selectrequest.Country{
					Code: payload.Location.Country,
				},
			},
		},
		Message: selectrequest.Message{
			Order: selectrequest.Order{
				Provider: selectrequest.Provider{
					ID: payload.ProviderID,
				},
				Fulfillments: []selectrequest.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
					},
				},
				Items: []selectrequest.Items{
					{
						ID:             payload.JobID,
					},
				},
			},
		},
	}
	return &req, nil
}
