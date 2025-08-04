package api

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/ChayanDass/beneficiary-manager/cmd/laas/docs"
	"github.com/ChayanDass/beneficiary-manager/pkg/middleware"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router Get the gin router with all the routes defined
//
//	@title						Beneficiary Manager API
//	@version					3.0
//	@description				Sochalarship  query over REST API.
//
//	@contact.name				chayan das
//	@contact.url				chayandass.com
//	@contact.email				daschayan8837@gmail.com
//
//	@license.name				GPL-2.0-only
//	@license.url				https://opensource.org/licenses/GPL-2.0
//
//	@BasePath					/api/v1
//	@securityDefinitions.basic	BasicAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your username and password for basic authentication

const (
	DEFAULT_PORT = "8080"
)

func Router() *gin.Engine {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = DEFAULT_PORT
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)

	oldSecurityScheme := regexp.MustCompile(`({\s*"ApiKeyAuth":\s*\[\]),\s*"{}":\s*\[\](\s*})`)
	docs.SwaggerInfo.SwaggerTemplate = oldSecurityScheme.ReplaceAllString(docs.SwaggerInfo.SwaggerTemplate, "$1$2, {}")

	// Initialize Gin router
	r := gin.Default()

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	// Handle invalid routes
	r.NoRoute(HandleInvalidUrl)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	api := r.Group("/api/v1")
	{
		// schemes routes group
		schemes := api.Group("/schemes")
		{
			schemes.GET("", GetSchemes)
			schemeByID := schemes.Group("/:id")
			{
				schemeByID.GET("", GetSchemeByID)
				schemeByID.GET("/status", GetSchemeStatus)
				schemeByID.GET("/applications", GetSchemeApplications)
			}
		}

		// User Routes
		user := api.Group("/users")
		{
			user.POST("/", CreateUser) // Create a new user
		}

		// Application Routes
		application := api.Group("/applications")
		application.Use(middleware.BasicAuth())
		{
			application.POST("/", SubmitApplication)                       // Submit application
			application.GET("/", GetApplications)                          // Get application status
			application.POST("/withdraw-application", WithdrawApplication) // Submit application without user ID
			application.POST("/init-application", InitApplication)         // Initialize application
			application.PATCH("/:id", ModifyApplication)                   // Update application
			application.GET("/status/:id", GetApplicationStatus)           // Get application by ID

		}

	}
	return r
}

func HandleInvalidUrl(c *gin.Context) {
	er := models.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "No such path exists, please check the URL",
		Error:   "invalid path",
	}
	c.JSON(http.StatusNotFound, er)
}
