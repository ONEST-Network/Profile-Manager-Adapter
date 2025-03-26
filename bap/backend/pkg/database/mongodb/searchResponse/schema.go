package searchresponse

import (
	cancelresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/response"
	confirmresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/response"
	initresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/response"
	searchresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response"
	selectresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/response"
	statusresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/response"
)

type SearchJobResponse struct {
	ID                 string                            `bson:"id"`
	TransactionID      string                            `bson:"transaction_id"`
	JobsResponse       []searchresponse.SearchResponse   `bson:"jobs_response"`
	SelectJobResponse  []selectresponse.SelectResponse   `bson:"select_job_response"`
	InitJobResponse    []initresponse.InitResponse       `bson:"init_job_response"`
	ConfirmJobResponse []confirmresponse.ConfirmResponse `bson:"confirm_job_response"`
	StatusJobResponse  []statusresponse.StatusResponse   `bson:"status_job_response"`
	CancelJobResponse  []cancelresponse.CancelResponse   `bson:"cancel_job_response"`
}
