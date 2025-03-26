package jobrecommender

import (
    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
)

// JobRecommendationRequestWithJobs extends the basic request to include jobs for filtering
type JobRecommendationRequestWithJobs struct {
    workerProfile.WorkerProfile
    Jobs []job.Job `json:"jobs,omitempty"` // Jobs to be filtered and recommended
}

// JobData represents a single job recommendation
type JobData struct {
    JobID       string `json:"job_id"`
    Title       string `json:"title"`
    Positions   string `json:"positions"`
    Skills      string `json:"skills"`
    Location    string `json:"location"`
    Salary      string `json:"salary"`
    Description string `json:"description"`
}

// JobRecommendationResponse represents the response from the recommendation service
type JobRecommendationResponse struct {
    Prompt         string    `json:"prompt"`
    ResultsFound   int       `json:"results_found"`
    Recommendations []JobData `json:"recommendations"`
}