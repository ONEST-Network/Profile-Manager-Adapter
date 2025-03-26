package handlers

import (
    "net/http"

    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/job-recommender"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// JobRecommendationHandler handles job recommendation related endpoints
type JobRecommendationHandler struct {
    clients *clients.Clients
}

// NewJobRecommendationHandler creates a new job recommendation handler
func NewJobRecommendationHandler(clients *clients.Clients) *JobRecommendationHandler {
    return &JobRecommendationHandler{
        clients: clients,
    }
}

// GetJobRecommendations handles the GET /jobs/recommend endpoint
// @Summary Get job recommendations
// @Description Get personalized job recommendations based on user profile
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body models.JobRecommendationRequest true "Job Recommendation Request"
// @Success 200 {object} models.JobRecommendationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /jobs/recommend [post]
func (h *JobRecommendationHandler) GetJobRecommendations(c *gin.Context) {
    var req jobrecommender.JobRecommendationRequestWithJobs
    if err := c.ShouldBindJSON(&req); err != nil {
        logrus.Errorf("Failed to bind request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validate request
    if req.Name == "" || req.Age == 0 || req.Gender == "" || req.PreferredLanguages == nil || 
       len(req.PreferredJobRoles) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
        return
    }

    // Apply filters based on user preferences
    filteredJobs := getSearchFilter(&req.WorkerProfile)
    logrus.Infof("Filtered down to %d relevant jobs based on user preferences", len(filteredJobs))
    

	jobs, err := h.clients.JobClient.ListJobs(filteredJobs)
	if err != nil {
		logrus.Errorf("Failed to list filtered jobs, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Jobs = jobs

    recommendations, err := h.clients.RecommendationClient.GetRecommendedJobs(&req)
    if err != nil {
        logrus.Errorf("Failed to get job recommendations: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, recommendations)
}

func getSearchFilter(payload *workerProfile.WorkerProfile) bson.D {
	var (
		locations = payload.Location
		query     = bson.D{}
	)

	// Use a zero value check for struct instead of nil comparison
	if locations != (workerProfile.Location{}) {
		var location = locations

		if location.City != "" {
			query = append(query, bson.E{Key: "location.city", Value: location.City})
		}
		if location.State != "" {
			query = append(query, bson.E{Key: "location.state", Value: location.State})
		}
		if location.AreaCode != "" {
			query = append(query, bson.E{Key: "location.area_code", Value: location.AreaCode})
			query = append(query, bson.E{Key: "$or", Value: bson.A{bson.D{{
				Key: "location.address",
				Value: bson.D{
					{Key: "$regex", Value: location.AreaCode},
					{Key: "$options", Value: "i"},
				},
			}}}})
		}
		if location.Coordinates.Longitute != 0 && location.Coordinates.Latitude != 0 {
			query = append(query, bson.E{
				Key: "address.coordinates", Value: bson.D{
					{Key: "$nearSphere", Value: bson.D{
						{Key: "$geometry", Value: bson.D{
							{Key: "type", Value: "Point"},
							{
								Key: "coordinates",
								Value: bson.A{
									location.Coordinates.Longitute,
									location.Coordinates.Latitude,
								},
							},
						}},
						{Key: "$maxDistance", Value: 5000}, // 5 km in meters
					}},
				},
			})
		}
	}

	return query
}
