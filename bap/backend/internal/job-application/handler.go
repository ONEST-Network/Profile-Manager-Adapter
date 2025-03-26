package jobapplication

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
)

type Interface interface {
}

type JobApplication struct {
	clients *clients.Clients
}

func NewJobApplication(clients *clients.Clients) Interface {
	return &JobApplication{
		clients: clients,
	}
}
