package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func WorkerProfileRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/add", handlers.AddWorkerProfile(clients))
}
