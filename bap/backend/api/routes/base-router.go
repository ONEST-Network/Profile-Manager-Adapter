package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
)

func BaseRouter(router *gin.RouterGroup, clients *clients.Clients) {
	// general routers
	router.GET("/status", handlers.StatusHandler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
