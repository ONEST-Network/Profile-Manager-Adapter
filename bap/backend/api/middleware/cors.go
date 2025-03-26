package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
)

const (
	AllowedOrigin      string = "Access-Control-Allow-Origin"
	AllowedMethods     string = "Access-Control-Allow-Methods"
	AllowedHeaders     string = "Access-Control-Allow-Headers"
	AllowedCredentials string = "Access-Control-Allow-Credentials"
)

func ValidateCors() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Supported origins
		allowedOrigins := config.Config.AllowedOrigins

		// Fetching the origin from the request context
		origin := c.GetHeader("Origin")

		// Validating the origin with the supported origin regex
		if origin != "" {
			validOrigin := false
			for _, allowedOrigin := range allowedOrigins {
				match, err := regexp.MatchString(allowedOrigin, origin)
				if err == nil && match {
					validOrigin = true
					c.Writer.Header().Set(AllowedOrigin, origin)
					break
				}
			}

			// if none of the origins are valid, returning error
			if !validOrigin {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Invalid origin",
				})
				c.Abort()
				return
			}
		}

		// Setting allowed methods
		c.Writer.Header().Set(AllowedMethods, strings.Join([]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		}, ","))

		// Setting allowed headers
		c.Writer.Header().Set(AllowedHeaders, "*")

		// Setting allowed credentials
		c.Writer.Header().Set(AllowedCredentials, "true")

		c.Next()
	}
}
