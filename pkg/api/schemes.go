package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/ChayanDass/beneficiary-manager/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @GetSchemes retrieves a list of schemes with pagination and filtering
// @Summary		Get Schemes with pagination and filtering
// @Description	Fetch available schemes with pagination and filtering
// @Tags			scheme
// @Accept			json
// @Produce		json
// @Param			page					query		int		false	"Page number"						default(1)
// @Param			limit					query		int		false	"Number of items per page"			default(10)
// @Param			name					query		string	false	"Scheme name"						example(Scholar Scheme)
// @Param			status					query		string	false	"Scheme status"						example(upcoming)
// @Param			min_amount				query		number	false	"Minimum amount"					example(1000)
// @Param			max_amount				query		number	false	"Maximum amount"					example(5000)
// @Param			start_after				query		string	false	"Start after date (RFC3339 format)"	example(2023-01-01T00:00:00Z)
// @Param			end_before				query		string	false	"End before date (RFC3339 format)"	example(2023-12-31T23:59:59Z)
// @Param			gender					query		string	false	"Eligible gender"
// @Param			academic_qualification	query		string	false	"Academic qualification"
// @Param			income_limit			query		number	false	"Income limit"
// @Param			category				query		string	false	"Category"
// @Success		200						{object}	models.SchemeResponse
// @Failure		400						{object}	models.ErrorResponse
// @Failure		500						{object}	models.ErrorResponse
// @Router			/schemes [get]
func GetSchemes(c *gin.Context) {
	pagination, offset := utils.GetPagination(c)

	var filter models.SchemeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	var schemes []models.Scheme
	query := db.DB.
		Preload("Eligibility").
		Preload("Eligibility.DocumentMappings").
		Preload("Eligibility.DocumentMappings.Document").
		Model(&models.Scheme{}).
		Joins("JOIN eligibilities ON eligibilities.id = schemes.eligibility_id")

	// Apply filters
	query = utils.ApplySchemeFilters(query, filter)

	// Count total
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch total count",
			Error:   err.Error(),
		})
		return
	}

	meta := utils.BuildPaginationMeta(c, pagination, totalCount)

	// Paginate and fetch
	if err := query.Offset(int(offset)).Limit(int(pagination.Limit)).Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch schemes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SchemeResponse{
		Data:    schemes,
		Code:    http.StatusOK,
		Message: "Schemes fetched successfully",
		Meta:    meta,
	})
}

//	@Summary		Get Scheme By ID
//	@Description	Get a specific scheme by its ID
//	@Tags			scheme
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Scheme ID"
//	@Success		200	{object}	models.Scheme
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/schemes/{id} [get]
//
// GetSchemeByID retrieves a scheme by its ID
func GetSchemeByID(c *gin.Context) {
	id := c.Param("id")
	var scheme models.Scheme

	if err := db.DB.
		Preload("Eligibility").
		Preload("Eligibility.DocumentMappings").
		Preload("Eligibility.DocumentMappings.Document").
		First(&scheme, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch scheme",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Scheme retrieved successfully",
		Data:    scheme,
	})
}

// @Summary		Get Scheme Status
// @Description	Get the current status of a scheme by its ID
// @Tags			scheme
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Scheme ID"
// @Success		200	{object}	models.SuccessResponse
// @Failure		404	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router		/schemes/{id}/status [get]
func GetSchemeStatus(c *gin.Context) {
	// Extract parameters from URL
	id := c.Param("id")
	var scheme models.Scheme

	err := db.DB.Where("id = ?", id).First(&scheme).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve scheme",
			Error:   err.Error(),
		})
		return
	}

	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application status fetched successfully",
		Data:    fmt.Sprintf("application status is %s", scheme.Status),
	}
	c.JSON(http.StatusOK, res)
}

// GetSchemeApplications retrieves all applications for a specific scheme.
//
// @Summary     Get all the applications for a scheme
// @Description Get all the applications for a scheme by its ID
// @Tags        scheme
// @Accept      json
// @Produce     json
// @Param       id      path     string  true   "Scheme ID"
// @Param       status  query    string  false  "Application status"  example(submitted)
// @Success     200     {object} models.SuccessResponse
// @Failure     400     {object} models.ErrorResponse
// @Router      /schemes/{id}/applications [get]
func GetSchemeApplications(c *gin.Context) {
	schemeID := c.Param("id")
	if schemeID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Scheme ID is required",
			Error:   "Missing 'scheme_id' path parameter",
		})
		return
	}

	// Default status = "submitted" unless specified otherwise
	status := c.DefaultQuery("status", "submitted")

	var applications []models.Application
	err := db.DB.Where("scheme_id = ? AND status = ?", schemeID, status).Find(&applications).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve applications",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Applications fetched successfully",
		Data:    applications,
	})
}
