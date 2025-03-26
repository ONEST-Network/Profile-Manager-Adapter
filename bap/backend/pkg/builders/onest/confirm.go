package onest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
	"github.com/google/uuid"
)

func BuildBPPConfirmJobRequest(payload confirmrequest.SeekerConfirmPayload, transactionId, messageId, bppId, bppuri string, worker *workerProfile.WorkerProfile) (*confirmrequest.ConfirmRequest, error) {
	languages, err := getConfirmLanguages(worker)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker languages: %v", err)
	}
	req := confirmrequest.ConfirmRequest{
		Context: confirmrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "confirm",
			Version:       "2.0.0",
			TransactionID: transactionId,
			MessageID:     messageId,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: confirmrequest.Location{
				City: confirmrequest.City{
					Code: payload.Location.City,
				},
				Country: confirmrequest.Country{
					Code: payload.Location.Country,
				},
			},
		},
		Message: confirmrequest.Message{
			Order: confirmrequest.Order{
				ID: uuid.New().String(),
				Provider: confirmrequest.Provider{
					ID: payload.ProviderID,
				},
				Fulfillments: []confirmrequest.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
						Customer: confirmrequest.Customer{
							Person: confirmrequest.Person{
								Name: worker.Name,
								Gender: string(worker.Gender),
								Age: strconv.Itoa(worker.Age),
								Skills:    getConfirmSkills(worker),
								Languages: languages,
								Creds:     getConfirmCreds(worker),
							},
							Contact: confirmrequest.Contact{
								Email: worker.Email,
								Phone: worker.Phone,
							},
						},
					},
				},
				Items: []confirmrequest.Items{
					{
						ID:             payload.JobID,
					},
				},
			},
		},
	}
	return &req, nil
}

func getConfirmSkills(worker *workerProfile.WorkerProfile) []confirmrequest.Skills {
	var skills []confirmrequest.Skills

	for _, skill := range worker.Skills {
		skills = append(skills, confirmrequest.Skills{
			Name: skill,
		})
	}

	return skills
}

func getConfirmLanguages(worker *workerProfile.WorkerProfile) ([]confirmrequest.Languages, error) {
	var languages []confirmrequest.Languages

	for _, language := range worker.PreferredLanguages {
		languageCode, err := utils.GetLanguageCode(string(language))
		if err != nil {
			return nil, err
		}
		languages = append(languages, confirmrequest.Languages{
			Name: languageCode,
		})
	}

	return languages, nil
}

func getConfirmCreds(worker *workerProfile.WorkerProfile) []confirmrequest.Creds {
	var creds []confirmrequest.Creds

	for i, cred := range worker.Credentials {
		creds = append(creds, confirmrequest.Creds{
			ID: "D"+strconv.Itoa(i+1),
			Descriptor: confirmrequest.CredsDescriptor{
				Name:      cred.Name,
				ShortDesc: cred.ShortDesc,
				LongDesc:  cred.LongDesc,
			},
			URL:  cred.URL,
			Type: cred.Type,
		})
	}

	return creds
}
