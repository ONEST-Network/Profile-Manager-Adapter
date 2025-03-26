package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func BecknRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/on_search", handlers.SearchJobs(clients))
	router.POST("/on_select", handlers.SendJobFulfillment(clients))
	router.POST("/on_init", handlers.InitializeJobApplication(clients))
	router.POST("/on_confirm", handlers.ConfirmJobApplication(clients))
	router.POST("/on_status", handlers.JobApplicationStatus(clients))
	router.POST("/on_cancel", handlers.WithdrawJobApplication(clients))
}

func BPPOnestRoutes(router *gin.RouterGroup, handler *handlers.OnestBPPHandler) {
    router.POST("/search", handler.Search())
    router.POST("/select", handler.Select())
    router.POST("/init", handler.Init())
    router.POST("/confirm", handler.Confirm())
    router.POST("/status", handler.Status())
    router.POST("/cancel", handler.Cancel())
}
