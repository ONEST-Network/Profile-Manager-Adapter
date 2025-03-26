package onest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
)

func BuildBPPInitJobRequest(payload initrequest.SeekerInitPayload, transactionId, messageId, bppId, bppuri string, worker *workerProfile.WorkerProfile) (*initrequest.InitRequest, error) {
	languages, err := getInitLanguages(worker)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker languages: %v", err)
	}
	req := initrequest.InitRequest{
		Context: initrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "init",
			Version:       "2.0.0",
			TransactionID: transactionId,
			MessageID:     messageId,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: initrequest.Location{
				City: initrequest.City{
					Code: payload.Location.City,
				},
				Country: initrequest.Country{
					Code: payload.Location.Country,
				},
			},
		},
		Message: initrequest.Message{
			Order: initrequest.Order{
				Provider: initrequest.Provider{
					ID: payload.ProviderID,
				},
				Fulfillments: []initrequest.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
						Customer: initrequest.Customer{
							Person: initrequest.Person{
								Name: worker.Name,
								Gender: string(worker.Gender),
								Age: strconv.Itoa(worker.Age),
								Skills:    getInitSkills(worker),
								Languages: languages,
								Creds:     getInitCreds(worker),
								Tags: []initrequest.Tags{
									{
										Descriptor: initrequest.TagsDescriptor{
											Code: "WORK_EXPERIENCE",
										},
										List: []initrequest.List{
											{
												Code: "TOTAL_EXPERIENCE",
												Value: "P" + strconv.Itoa(worker.Experience) + "Y",
											},
										},
									},
								},
							},
							Contact: initrequest.Contact{
								Email: worker.Email,
								Phone: worker.Phone,
							},
						},
					},
				},
				Items: []initrequest.Items{
					{
						ID:             payload.JobID,
					},
				},
			},
		},
	}
	return &req, nil
}

func getInitSkills(worker *workerProfile.WorkerProfile) []initrequest.Skills {
	var skills []initrequest.Skills

	for _, skill := range worker.Skills {
		skills = append(skills, initrequest.Skills{
			Name: skill,
		})
	}

	return skills
}

func getInitLanguages(worker *workerProfile.WorkerProfile) ([]initrequest.Languages, error) {
	var languages []initrequest.Languages

	for _, language := range worker.PreferredLanguages {
		languageCode, err := utils.GetLanguageCode(string(language))
		if err != nil {
			return nil, err
		}
		languages = append(languages, initrequest.Languages{
			Name: languageCode,
		})
	}

	return languages, nil
}

func getInitCreds(worker *workerProfile.WorkerProfile) []initrequest.Creds {
	var creds []initrequest.Creds

	for i, cred := range worker.Credentials {
		creds = append(creds, initrequest.Creds{
			ID: "D"+strconv.Itoa(i+1),
			Descriptor: initrequest.CredsDescriptor{
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
