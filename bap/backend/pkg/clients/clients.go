package clients

import (
	apiclient "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/api-client"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/cache/redis"
	dbJob "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
	dbWorker "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	dbSearchResponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/searchResponse"
	dbInitJobApplication "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/init-job-application"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/job-recommender"
)

type Clients struct {
	ApiClient           apiclient.Interface
	RedisClient 	   *redis.RedisClient
	WorkerProfileClient *dbWorker.Dao
	SearchReponseClient *dbSearchResponse.Dao
	JobClient           *dbJob.Dao
	InitJobApplicationClient *dbInitJobApplication.Dao
	RecommendationClient *jobrecommender.JobRecommendationClient
}

func NewClients(jobClient *dbJob.Dao, workerProfileClient *dbWorker.Dao, searchReponseClient *dbSearchResponse.Dao, redisClient *redis.RedisClient) *Clients {
	return &Clients{
		RedisClient: redisClient,
		ApiClient:           apiclient.NewAPIClient(),
		WorkerProfileClient: workerProfileClient,
		JobClient:                jobClient,
		RecommendationClient: jobrecommender.NewJobRecommendationClient(),
		SearchReponseClient: searchReponseClient,
	}
}
