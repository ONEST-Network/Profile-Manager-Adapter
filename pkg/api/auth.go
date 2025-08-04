package api

import (
	"net/http"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
)

// CreateUser  user create to test all the other endpoints
// @Summary      Create a new user
// @Description  Registers a new user with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreate  true  "User create payload"
// @Success      201   {object}  models.User
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var req models.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
