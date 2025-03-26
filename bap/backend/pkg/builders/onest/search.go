package onest

import (
	"time"

	"github.com/google/uuid"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
)

func BuildBPPSearchJobsRequest(payload searchrequest.SeekerSearchPayload) (*searchrequest.SearchRequest, string, error) {
	transaction_id := uuid.New().String()
	req := searchrequest.SearchRequest{
		Context: searchrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "search",
			Version:       "2.0.0",
			TransactionID: transaction_id,
            MessageID:     uuid.New().String(),
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: searchrequest.Location{
				City: searchrequest.City{
					Code: payload.Location.City,
				},
				Country: searchrequest.Country{
					Code: payload.Location.Country,
				},
			},
		},
		Message: searchrequest.Message{
			Intent: searchrequest.Intent{
				Item: searchrequest.Item{
					Descriptor: searchrequest.ItemDescriptor{
						Name: payload.Role,
					},
				},
				Provider: searchrequest.Provider{
					Descriptor: searchrequest.ProviderDescriptor{
						Name: payload.Provider,
					},
					Locations: []searchrequest.ProviderLocations{
						{
							City: searchrequest.ProviderCity {
								Code: payload.Location.City,
							},
							State: searchrequest.ProviderState {
								Code: payload.Location.State,
							},
							AreaCode: searchrequest.ProviderAreaCode {
								Code: payload.Location.AreaCode,
							},
							Coordinates: searchrequest.Coordinates{
								Latitude:  payload.Location.Coordinates.Latitude,
								Longitude: payload.Location.Coordinates.Longitude,
							},
						},
					},
				},
				Tags: []searchrequest.Tags{
					{
						Descriptor: searchrequest.TagsDescriptor{
							Code: "JOB_DETAILS",
						},
						List: []searchrequest.List{
							{
								Descriptor: searchrequest.TagsDescriptor{
									Code: "INDUSTRY_TYPE",
								},
								Value: payload.Category,
							},
							{
								Descriptor: searchrequest.TagsDescriptor{
									Code: "JOB_TYPE",
								},
								Value: payload.EmploymentType,
							},
						},
					},
				},
			},
		},
	}
	return &req, transaction_id, nil
}
