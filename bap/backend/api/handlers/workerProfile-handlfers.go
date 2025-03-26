package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/workerProfile"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	workerProfilePayload "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/workerProfile"
	"github.com/gin-gonic/gin"
)

func AddWorkerProfile(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload workerProfilePayload.AddWorkerProfileRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := workerProfile.NewWorkerProfile(clients).AddWorkerProfile(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}