package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/service"
	builders "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/builders/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	dbSearchResponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/searchResponse"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
)

type OnestBPPHandler struct {
	onestService *service.OnestBPPService
}

func NewOnestBPPHandler(onestService *service.OnestBPPService) *OnestBPPHandler {
	return &OnestBPPHandler{
		onestService: onestService,
	}
}

func SearchJobs(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SearchJobsAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		go onest.SearchJobs(payload)
	}
}

func SendJobFulfillment(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SendJobFulfillmentAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.SendJobFulfillment(payload)
	}
}

func InitializeJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.InitializeJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.InitializeJobApplication(payload)
	}
}

func ConfirmJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.ConfirmJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.ConfirmJobApplication(payload)
	}
}

func JobApplicationStatus(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.JobApplicationStatusAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.JobApplicationStatus(payload)
	}
}

func WithdrawJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.WithdrawJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.WithdrawJobApplication(payload)
	}
}

/**
BPP APIs
**/

// @Summary	Search jobs
// @Description	Search jobs
// @Tags Worker
// @Accept		json
// @Produce		json
// @Param request body searchrequest.SeekerSearchPayload true "request body"
// @Success 200 {object} searchresponse.SearchResponse
// @Failure 500 {object} string
// @Router	/search	[post]

func (h *OnestBPPHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload searchrequest.SeekerSearchPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        // Store transaction ID and message ID in worker profile
        if payload.WorkerID != "" {
            parsedRequest, transaction_id, err := builders.BuildBPPSearchJobsRequest(payload)
            if err != nil {
                logrus.Errorf("Failed to parse search job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Create a new document in the SearchResponseClient collection
            searchJobResponse := dbSearchResponse.SearchJobResponse {
                ID: transaction_id,
                TransactionID: transaction_id,
            }

            // Cache the initial search response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                transaction_id,
                "searching",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache search response in Redis for transaction_id %s: %v", transaction_id, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            // Insert the document into the collection
            err = h.onestService.Clients.SearchReponseClient.CreateSearchJobResponse(&searchJobResponse)
            if err != nil {
                logrus.Errorf("Failed to create search response document: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create search response"})
                return
            }

            logrus.Infof("Created search response document with transaction_id: %s", transaction_id)

            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
            updateFields := bson.D{{Key: "$set", Value: bson.D{
                {Key: "last_transaction_id", Value: parsedRequest.Context.TransactionID},
                {Key: "message_id", Value: parsedRequest.Context.MessageID},
            }}}
            
            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                logrus.Errorf("Failed to update worker profile with transaction ID %s: %v",parsedRequest.Context.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            logrus.Printf("Parsed Seeker Search Request: %+v\n", parsedRequest)
            _, err = h.onestService.Search(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send search request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            
            // Wait for a matching document
            logrus.Infof("Waiting for search job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(transaction_id, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", transaction_id, err)
                        continue
                    }
                    if status == "got response" {
                        // Get the search response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            transaction_id,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get search response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(transaction_id, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", transaction_id, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            c.JSON(http.StatusOK, searchResponse.JobsResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(transaction_id, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", transaction_id, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No jobs found"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(transaction_id, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", transaction_id, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Select() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload selectrequest.SeekerSelectPayload
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            workerTransactionId := worker.TransactionID
            workerMessageId := worker.MessageID

            activeJobApplication, exists := worker.ActiveJobApplications[payload.JobID]
            if !exists {
                // This means the job ID doesn't exist in active applications
                updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                updateFields := bson.D{{Key: "$set", Value: bson.D{
                    {Key: "active_job_applications." + payload.JobID, Value: bson.M{
                        "transaction_id": workerTransactionId,
                        "bpp_id": payload.BppID,
                        "bpp_uri": payload.BppURI,
                    }},
                }}}
                
                if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                    logrus.Errorf("Failed to update worker profile with active job applications for transaction ID %s: %v", workerTransactionId, err)
                    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                    return
                }
            } else {
                // This means the job ID already exists in active applications
                if activeJobApplication.LastRequestExecuted == "select" {
                    // This means the job ID is already in the process of being selected
                    c.JSON(http.StatusConflict, gin.H{"error": "Job has already been selected"})
                    return
                }
            }

            parsedRequest, err := builders.BuildBPPSelectJobRequest(payload, workerTransactionId, workerMessageId, payload.BppID, payload.BppURI)
            if err != nil {
                logrus.Errorf("Failed to parse select job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }


            // Cache the initial select response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                workerTransactionId,
                "selecting",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache select response in Redis for transaction_id %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            }

            logrus.Printf("Parsed Seeker Select Request: %+v\n", parsedRequest)


            _, err = h.onestService.Select(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send select request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            // Wait for a matching document
            logrus.Infof("Waiting for select job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(worker.TransactionID, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        continue
                    }
                    if status == "selected job" {
                        // Get the search response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            worker.TransactionID,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get select response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                            updateFields := bson.D{{Key: "$set", Value: bson.D{
                                {Key: "active_job_applications." + payload.JobID + ".last_request_executed", Value: "select"},
                            }}}

                            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                                logrus.Errorf("Failed to update last request executed status to select: %v", err)
                                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                                return
                            }
                            c.JSON(http.StatusOK, searchResponse.SelectJobResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No job selected"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Init() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload initrequest.SeekerInitPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            if worker.ActiveJobApplications[payload.JobID].LastRequestExecuted == "init" {
                // This means the job application for this job ID is already submitted
                c.JSON(http.StatusConflict, gin.H{"error": "Job application has already been submitted"})
                return
            }
            parsedRequest, err := builders.BuildBPPInitJobRequest(payload, worker.TransactionID, worker.MessageID, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse init job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            // Cache the initial init response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                worker.TransactionID,
                "filling job application",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache init response in Redis for transaction_id %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            logrus.Printf("Parsed Seeker Init Request: %+v\n", parsedRequest)
            _, err = h.onestService.Init(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send init request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            // Wait for a matching document
            logrus.Infof("Waiting for init job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(worker.TransactionID, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        continue
                    }
                    if status == "init job application" {
                        // Get the search response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            worker.TransactionID,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get init response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                            updateFields := bson.D{{Key: "$set", Value: bson.D{
                                {Key: "active_job_applications." + payload.JobID + ".last_request_executed", Value: "init"},
                            }}}

                            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                                logrus.Errorf("Failed to update last request executed status to init: %v", err)
                                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                                return
                            }
                            c.JSON(http.StatusOK, searchResponse.InitJobResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No job application found"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Confirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload confirmrequest.SeekerConfirmPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            if worker.ActiveJobApplications[payload.JobID].LastRequestExecuted == "confirm" {
                // This means the job application for this job ID has already been confirmed
                c.JSON(http.StatusConflict, gin.H{"error": "Job application has already been confirmed"})
                return
            }
            parsedRequest, err := builders.BuildBPPConfirmJobRequest(payload, worker.TransactionID, worker.MessageID, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse confirm job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
            updateFields := bson.D{{Key: "$set", Value: bson.D{
                {Key: "active_job_applications." + payload.JobID + ".application_id", Value: parsedRequest.Message.Order.ID},
            }}}

            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                logrus.Errorf("Failed to update application id in the user's active job applications: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Cache the initial confirm response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                worker.TransactionID,
                "confirming job application",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache confirm response in Redis for transaction_id %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            logrus.Printf("Parsed Seeker Confirm Request: %+v\n", parsedRequest)

            _, err = h.onestService.Confirm(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send confirm request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Wait for a matching document
            logrus.Infof("Waiting for confirm job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(worker.TransactionID, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        continue
                    }
                    if status == "confirm job application" {
                        // Get the search response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            worker.TransactionID,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get confirm response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                            updateFields := bson.D{{Key: "$set", Value: bson.D{
                                {Key: "active_job_applications." + payload.JobID + ".last_request_executed", Value: "confirm"},
                            }}}

                            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                                logrus.Errorf("Failed to update last request executed status to confirm: %v", err)
                                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                                return
                            }
                            c.JSON(http.StatusOK, searchResponse.ConfirmJobResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No job application found to be confirmed"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
            
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Status() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload statusrequest.SeekerStatusPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPStatusJobRequest(payload, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse status job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Cache the initial status response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                worker.TransactionID,
                "checking job application status",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache status response in Redis for transaction_id %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            }
            logrus.Printf("Parsed Seeker Status Request: %+v\n", parsedRequest)
            _, err = h.onestService.Status(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send status request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Wait for a matching document
            logrus.Infof("Waiting for status job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(worker.TransactionID, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        continue
                    }
                    if status == "status job application" {
                        // Get the search response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            worker.TransactionID,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get status response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            c.JSON(http.StatusOK, searchResponse.StatusJobResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No job application status found"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload cancelrequest.SeekerCancelPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            if worker.ActiveJobApplications[payload.ApplicationID].LastRequestExecuted == "cancel" {
                // This means the job application for this job ID has already been cancelled
                c.JSON(http.StatusConflict, gin.H{"error": "Job application has already been cancelled"})
                return
            }
            parsedRequest, err := builders.BuildBPPCancelJobRequest(payload, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse cancel job request, %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Cache the initial init response in Redis
            err = h.onestService.Clients.RedisClient.Set(
                worker.TransactionID,
                "cancelling job application",
                0,
            )
            if err != nil {
                logrus.Errorf("Failed to cache cancel response in Redis for transaction_id %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            }
            logrus.Printf("Parsed Seeker Cancel Request: %+v\n", parsedRequest)
            _, err = h.onestService.Cancel(c.Request.Context(), parsedRequest)
            if err != nil {
                logrus.Errorf("Failed to send cancel request to ONEST Network: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Wait for a matching document
            logrus.Infof("Waiting for cancel job response for the transaction ID: %s", parsedRequest.Context.TransactionID)
            ticker := time.NewTicker(1 * time.Second)
            defer ticker.Stop()

            pollTimeout := time.After(30 * time.Second) // Set a timeout for polling
            for {
                select {
                case <-ticker.C:
                    // Check Redis for status update
                    var status string
                    err = h.onestService.Clients.RedisClient.Get(worker.TransactionID, &status)
                    if err != nil {
                        logrus.Errorf("Failed to get Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        continue
                    }
                    if status == "cancel job application" {
                        // Get the cancel response from MongoDB
                        searchResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(
                            worker.TransactionID,
                        )
                        if err != nil {
                            logrus.Errorf("Failed to get cancel response from MongoDB: %v", err)
                            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
            
                        // Return the jobs response to the client
                        if searchResponse != nil && len(searchResponse.JobsResponse) > 0 {
                            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                            updateFields := bson.D{{Key: "$set", Value: bson.D{
                                {Key: "active_job_applications." + payload.JobID + ".last_request_executed", Value: "cancel"},
                            }}}

                            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                                logrus.Errorf("Failed to update last request executed status to cancel: %v", err)
                                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                                return
                            }
                            c.JSON(http.StatusOK, searchResponse.CancelJobResponse[0])
                            // Reset Redis status to prevent duplicate processing
                            if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                                logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                            }
                            return
                        }
                        
                        c.JSON(http.StatusNotFound, gin.H{"error": "No job application found to be cancelled"})
                        // Reset Redis status to prevent duplicate processing
                        if err := h.onestService.Clients.RedisClient.Set(worker.TransactionID, "", 0); err != nil {
                            logrus.Warnf("Failed to reset Redis status for transaction_id %s: %v", worker.TransactionID, err)
                        }
                        return
                    }
                case <-pollTimeout:
                    c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                    return
                }
            }
            
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}
