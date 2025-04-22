package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
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

func (h *OnestBPPHandler) Apply() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload selectrequest.SeekerSelectPayload
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Step 1: Select
        _, err := h.processSelect(payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"step": "failed to call select API", "error": err.Error()})
            return
        }

        // Step 2: Init
        _, err = h.processInit(payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"step": "failed to call init API", "error": err.Error()})
            return
        }

        // Step 3: Confirm
        confirmResponse, err := h.processConfirm(payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"step": "failed to call confirm API", "error": err.Error()})
            return
        }

        // Return combined response
        response := gin.H{
            "confirm": confirmResponse,
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestBPPHandler) processSelect(payload selectrequest.SeekerSelectPayload) (interface{}, error) {
    // Convert apply payload to select payload
    selectPayload := selectrequest.SeekerSelectPayload{
        WorkerID:   payload.WorkerID,
        JobID:      payload.JobID,
    }
    var response interface{}
    err := h.onestService.Clients.ApiClient.ApiCall(selectPayload, config.Config.BapUri + "/select", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to call select API: %v", err)
    }

    return response, nil
}

func (h *OnestBPPHandler) processInit(payload selectrequest.SeekerSelectPayload) (interface{}, error) {
    // Convert apply payload to init payload
    initPayload := initrequest.SeekerInitPayload{
        WorkerID:   payload.WorkerID,
        JobID:      payload.JobID,
    }

    var response interface{}
    err := h.onestService.Clients.ApiClient.ApiCall(initPayload, config.Config.BapUri + "/init", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to call init API: %v", err)
    }

    return response, nil
}

func (h *OnestBPPHandler) processConfirm(payload selectrequest.SeekerSelectPayload) (interface{}, error) {
    // Convert apply payload to confirm payload
    confirmPayload := confirmrequest.SeekerConfirmPayload{
        WorkerID:   payload.WorkerID,
        JobID:      payload.JobID,
    }

    var response interface{}
    err := h.onestService.Clients.ApiClient.ApiCall(confirmPayload, config.Config.BapUri + "/confirm", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to call confirm API: %v", err)
    }

    return response, nil
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
            bpp_id := ""
            bpp_uri := ""
            providerId := ""
            jobCityCode := ""
            jobCountryCode := ""
            // Check if job exists in search response
            searchJobResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(workerTransactionId)
            if err != nil {
                logrus.Errorf("Failed to get search job response for transaction ID %s: %v", workerTransactionId, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Verify if the job exists in the search response
            jobExists := false
            if searchJobResponse != nil && len(searchJobResponse.JobsResponse) > 0 {
                // Check first response's providers
                response := searchJobResponse.JobsResponse[0].Message.Catalog
                for _, provider := range response.Providers {
                    for _, item := range provider.Items {
                        if item.ID == payload.JobID {
                            jobExists = true
                            bpp_id = searchJobResponse.JobsResponse[0].Context.BppID
                            bpp_uri = searchJobResponse.JobsResponse[0].Context.BppURI
                            providerId = provider.ID
                            jobCityCode = item.Creator.City.Code
                            jobCountryCode = searchJobResponse.JobsResponse[0].Context.Location.Country.Code
                            break
                        }
                    }
                    if jobExists {
                        break
                    }
                }
            }
            if !jobExists {
                logrus.Errorf("Job ID %s not found in search response for transaction ID %s", payload.JobID, workerTransactionId)
                c.JSON(http.StatusNotFound, gin.H{"error": "Job not found in search results"})
                return
            }
            activeJobApplication, exists := worker.ActiveJobApplications[payload.JobID]
            if !exists {
                // This means the job ID doesn't exist in active applications
                updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
                updateFields := bson.D{{Key: "$set", Value: bson.D{
                    {Key: "active_job_applications." + payload.JobID, Value: bson.M{
                        "transaction_id": workerTransactionId,
                        "bpp_id": bpp_id,
                        "bpp_uri": bpp_uri,
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

            parsedRequest, err := builders.BuildBPPSelectJobRequest(payload, workerTransactionId, workerMessageId, bpp_id, bpp_uri, providerId, jobCityCode, jobCountryCode)
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

            bpp_id := ""
            bpp_uri := ""
            providerId := ""
            jobCityCode := ""
            jobCountryCode := ""
            // Check if job exists in select response
            selectJobResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(worker.TransactionID)
            if err != nil {
                logrus.Errorf("Failed to get select job response for transaction ID %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Verify if the job exists in the select response
            jobExists := false
            if selectJobResponse != nil && len(selectJobResponse.SelectJobResponse) > 0 {
                response := selectJobResponse.SelectJobResponse[0].Message.Order
                for _, item := range response.Items {
                    if item.ID == payload.JobID {
                        jobExists = true
                        bpp_id = selectJobResponse.SelectJobResponse[0].Context.BppID
                        bpp_uri = selectJobResponse.SelectJobResponse[0].Context.BppURI
                        providerId = response.Provider.ID
                        break
                    }
                    if jobExists {
                        break
                    }
                }
                jobExists = false 
                for _, provider := range selectJobResponse.JobsResponse[0].Message.Catalog.Providers {
                    for _, item := range provider.Items {
                        if item.ID == payload.JobID {
                            jobExists = true
                            jobCityCode = item.Creator.City.Code
                            jobCountryCode = selectJobResponse.JobsResponse[0].Context.Location.Country.Code
                            break
                        }
                    }
                    if jobExists {
                        break
                    }
                }
            }
            if !jobExists {
                logrus.Errorf("Job ID %s not found in select response for transaction ID %s", payload.JobID, worker.TransactionID)
                c.JSON(http.StatusNotFound, gin.H{"error": "Job not found in select results"})
                return
            }

            parsedRequest, err := builders.BuildBPPInitJobRequest(payload, worker.TransactionID, worker.MessageID, bpp_id, bpp_uri, worker, providerId, jobCityCode, jobCountryCode)
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

            bpp_id := ""
            bpp_uri := ""
            providerId := ""
            jobCityCode := ""
            jobCountryCode := ""
            // Check if job exists in init response
            initJobResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(worker.TransactionID)
            if err != nil {
                logrus.Errorf("Failed to get init job response for transaction ID %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Verify if the job exists in the init response
            jobExists := false
            if initJobResponse != nil && len(initJobResponse.InitJobResponse) > 0 {
                // Check first response's providers
                response := initJobResponse.InitJobResponse[0].Message.Order
                for _, item := range response.Items {
                    if item.ID == payload.JobID {
                        jobExists = true
                        bpp_id = initJobResponse.InitJobResponse[0].Context.BppID
                        bpp_uri = initJobResponse.InitJobResponse[0].Context.BppURI
                        providerId = response.Provider.ID
                        break
                    }
                    if jobExists {
                        break
                    }
                }
                jobExists = false 
                for _, provider := range initJobResponse.JobsResponse[0].Message.Catalog.Providers {
                    for _, item := range provider.Items {
                        if item.ID == payload.JobID {
                            jobExists = true
                            jobCityCode = item.Creator.City.Code
                            jobCountryCode = initJobResponse.JobsResponse[0].Context.Location.Country.Code
                            break
                        }
                    }
                    if jobExists {
                        break
                    }
                }
            }
            if !jobExists {
                logrus.Errorf("Job ID %s not found in init response for transaction ID %s", payload.JobID, worker.TransactionID)
                c.JSON(http.StatusNotFound, gin.H{"error": "Job not found in init results"})
                return
            }

            parsedRequest, err := builders.BuildBPPConfirmJobRequest(payload, worker.TransactionID, worker.MessageID, bpp_id, bpp_uri, worker, providerId, jobCityCode, jobCountryCode)
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

            bpp_id := ""
            bpp_uri := ""
            jobCityCode := ""
            jobCountryCode := ""
            jobId := ""
            // Check if job exists in confirm response
            confirmJobResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(worker.TransactionID)
            if err != nil {
                logrus.Errorf("Failed to get confirm job response for transaction ID %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Verify if the job exists in the confirm response
            applicationExists := false
            if confirmJobResponse != nil && len(confirmJobResponse.ConfirmJobResponse) > 0 {
                response := confirmJobResponse.ConfirmJobResponse[0].Message.Order
                if response.ID == payload.ApplicationID {
                    applicationExists = true
                    bpp_id = confirmJobResponse.ConfirmJobResponse[0].Context.BppID
                    bpp_uri = confirmJobResponse.ConfirmJobResponse[0].Context.BppURI
                    jobId = response.Items[0].ID
                }
                jobExists := false 
                for _, provider := range confirmJobResponse.JobsResponse[0].Message.Catalog.Providers {
                    for _, item := range provider.Items {
                        if item.ID == jobId {
                            jobExists = true
                            jobCityCode = item.Creator.City.Code
                            jobCountryCode = confirmJobResponse.JobsResponse[0].Context.Location.Country.Code
                            break
                        }
                    }
                    if jobExists {
                        break
                    }
                }
            }
            if !applicationExists {
                logrus.Errorf("Application ID %s not found in init response for transaction ID %s", payload.ApplicationID, worker.TransactionID)
                c.JSON(http.StatusNotFound, gin.H{"error": "Application not found in confirm results"})
                return
            }

            parsedRequest, err := builders.BuildBPPStatusJobRequest(payload, bpp_id, bpp_uri, worker, jobCityCode, jobCountryCode)
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

            bpp_id := ""
            bpp_uri := ""
            jobCityCode := ""
            jobCountryCode := ""
            // Check if job exists in status response
            statusJobResponse, err := h.onestService.Clients.SearchReponseClient.GetSearchJobResponse(worker.TransactionID)
            if err != nil {
                logrus.Errorf("Failed to get status job response for transaction ID %s: %v", worker.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            // Verify if the job exists in the status response
            applicationExists := false
            if statusJobResponse != nil && len(statusJobResponse.StatusJobResponse) > 0 {
                response := statusJobResponse.StatusJobResponse[0].Message.Order
                if response.ID == payload.ApplicationID {
                    applicationExists = true
                    bpp_id = statusJobResponse.StatusJobResponse[0].Context.BppID
                    bpp_uri = statusJobResponse.StatusJobResponse[0].Context.BppURI
                }
                jobExists := false 
                for _, provider := range statusJobResponse.JobsResponse[0].Message.Catalog.Providers {
                    for _, item := range provider.Items {
                        if item.ID == payload.JobID {
                            jobExists = true
                            jobCityCode = item.Creator.City.Code
                            jobCountryCode = statusJobResponse.JobsResponse[0].Context.Location.Country.Code
                            break
                        }
                    }
                    if jobExists {
                        break
                    }
                }
            }
            if !applicationExists {
                logrus.Errorf("Application ID %s not found in status response for transaction ID %s", payload.ApplicationID, worker.TransactionID)
                c.JSON(http.StatusNotFound, gin.H{"error": "Application not found in status results"})
                return
            }

            parsedRequest, err := builders.BuildBPPCancelJobRequest(payload, bpp_id, bpp_uri, worker, jobCityCode, jobCountryCode)
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
