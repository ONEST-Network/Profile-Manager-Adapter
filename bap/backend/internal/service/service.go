package service

import (
	"context"
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	cancelresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/response"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	confirmresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/response"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	initresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/response"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	searchresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	selectresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/response"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
	statusresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/response"
)

type OnestBPPService struct {
    Clients *clients.Clients
}

func NewOnestBPPService(clients *clients.Clients) *OnestBPPService {
    return &OnestBPPService{
        Clients: clients,
    }
}

func (s *OnestBPPService) Search(ctx context.Context, req *searchrequest.SearchRequest) (*searchresponse.SearchResponse, error) {
    // Make API call to BPP search endpoint
    var response searchresponse.SearchResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/search", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to search jobs: %w", err)
    }
    return &response, nil
}

func (s *OnestBPPService) Select(ctx context.Context, req *selectrequest.SelectRequest) (*selectresponse.SelectResponse, error) {
    // Make API call to BPP select endpoint
    var response selectresponse.SelectResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/select", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to select job: %w", err)
    }
    return &response, nil
}

func (s *OnestBPPService) Init(ctx context.Context, req *initrequest.InitRequest) (*initresponse.InitResponse, error) {
    // Make API call to BPP init endpoint
    var response initresponse.InitResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/init", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to initialize job application: %w", err)
    }
    return &response, nil
}

func (s *OnestBPPService) Confirm(ctx context.Context, req *confirmrequest.ConfirmRequest) (*confirmresponse.ConfirmResponse, error) {
    // Make API call to BPP confirm endpoint
    var response confirmresponse.ConfirmResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/confirm", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to confirm job application: %w", err)
    }
    return &response, nil
}

func (s *OnestBPPService) Status(ctx context.Context, req *statusrequest.StatusRequest) (*statusresponse.StatusResponse, error) {
    // Make API call to BPP status endpoint
    var response statusresponse.StatusResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/status", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to get job application status: %w", err)
    }
    return &response, nil
}

func (s *OnestBPPService) Cancel(ctx context.Context, req *cancelrequest.CancelRequest) (*cancelresponse.CancelResponse, error) {
    // Make API call to BPP cancel endpoint
    var response cancelresponse.CancelResponse
    err := s.Clients.ApiClient.ApiCall(req, config.Config.BppUri + "/cancel", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to cancel job application: %w", err)
    }
    return &response, nil
}

