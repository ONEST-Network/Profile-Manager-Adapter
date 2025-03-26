package main

import (
	"fmt"
	// "time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/log"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/proxy"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/server"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/service"
	// "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// init logrus logger
	log.InitLogger()

	// parse env variables
	if err := envconfig.Process("", &config.Config); err != nil {
		logrus.Fatalf("Failed to parse ENVs, %v", err)
	}

	// set proxy envs, if provided
	proxy.SetProxyENVs()

	// Initialize mongodb clients
	seekerClient, jobClient, searchJobResponse := server.InitMongoDB()

	redisClient := server.InitRedis()

	// Set up clients
	clients := clients.NewClients(jobClient, seekerClient, searchJobResponse, redisClient)

	// Set up the service
	onestBPPService := service.NewOnestBPPService(clients)

	// Set up the BPP proxy
	bppHandler := handlers.NewOnestBPPHandler(onestBPPService)

	// initialize the server
	server := server.SetupServer(clients, bppHandler)

	// // Initialize job sync from BPP side
	// jobSync := utils.NewJobSync(clients, 5*time.Minute)
    // jobSync.Start()
    // defer jobSync.Stop()

	// start the server
	fmt.Println("Starting server at port", config.Config.HTTPPort)
	if err := server.Run(fmt.Sprintf(":%s", config.Config.HTTPPort)); err != nil {
		logrus.Fatal(err)
	}
}
