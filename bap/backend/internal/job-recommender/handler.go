package jobrecommender

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
    "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/job-recommender"
    "github.com/sirupsen/logrus"
)

// JobRecommendationClient handles communication with the job recommendation service
type JobRecommendationClient struct {
    BaseURL    string
    HTTPClient *http.Client
}


// NewJobRecommendationClient creates a new client for the job recommendation service
func NewJobRecommendationClient() *JobRecommendationClient {
    return &JobRecommendationClient{
        BaseURL: config.Config.RecommendationServiceURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// GetRecommendedJobsWithFilter calls the recommendation service with user profile and job list
func (c *JobRecommendationClient) GetRecommendedJobs(req *jobrecommender.JobRecommendationRequestWithJobs) (*jobrecommender.JobRecommendationResponse, error) {
    url := fmt.Sprintf("%s/recommend_jobs", c.BaseURL)
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }
    
    return c.sendRequest(url, jsonData)
}

// sendRequest is a helper function to send HTTP requests
func (c *JobRecommendationClient) sendRequest(url string, jsonData []byte) (*jobrecommender.JobRecommendationResponse, error) {
    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    
    logrus.Infof("Sending job recommendation request to: %s", url)
    resp, err := c.HTTPClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("recommendation service returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
    }
    
    var recommendationResp jobrecommender.JobRecommendationResponse
    if err := json.NewDecoder(resp.Body).Decode(&recommendationResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }
    
    return &recommendationResp, nil
}

