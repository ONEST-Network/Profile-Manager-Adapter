// @title Beneficiary Manager API
// @version 1.0
// @description This is a backend service for managing schemes, users, and applications.
// @contact.name Chayan
// @host localhost:8080
// @BasePath /api
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ChayanDass/beneficiary-manager/pkg/api"
	"github.com/ChayanDass/beneficiary-manager/pkg/config"
	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

// main is the entry point of the application. It performs the following tasks:
// 1. Loads the database configuration using the `config.LoadDBConfig` function.
// 2. Loads environment variables from the `.env` file using `godotenv.Load`.
// 3. Parses command-line flags using `flag.Parse`.
// 4. Establishes a connection to the database using the `db.Connect` function.
// 5. Sets up the API router using the `api.Router` function.
// 6. Performs GORM auto-migration for various database models to ensure the database schema is up-to-date.
//   - Models include Application, User, StudentAcademicQualification, Scheme, Eligibility, Address,
//     StudentProfile, UploadDocument, DocumentsRequired, and EligibilityDocumentMap.
//
// 7. Seeds the database with default document types using the `DefaultDocumentsRequired` model.
// 8. Starts the HTTP server using the `r.Run` function and logs any errors encountered during execution.
func main() {

	cfg := config.LoadDBConfig()

	fmt.Printf("Connecting to DB at %s:%s...\n", cfg.Host, cfg.Port)
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	flag.Parse()
	db.Connect(&cfg.Host, &cfg.Port, &cfg.User, &cfg.DBName, &cfg.Password)
	r := api.Router()

	// Perform GORM auto-migration for the Application table
	if err := db.DB.AutoMigrate(&models.Application{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the User table
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the StudentAcademicQualification table
	if err := db.DB.AutoMigrate(&models.StudentAcademicQualification{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the Scheme table
	if err := db.DB.AutoMigrate(&models.Scheme{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the Eligibility table
	if err := db.DB.AutoMigrate(&models.Eligibility{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the Address table
	if err := db.DB.AutoMigrate(&models.Address{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the StudentProfile table
	if err := db.DB.AutoMigrate(&models.StudentProfile{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the UploadDocument table
	if err := db.DB.AutoMigrate(&models.UploadDocument{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the DocumentsRequired table
	if err := db.DB.AutoMigrate(&models.DocumentsRequired{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Perform GORM auto-migration for the EligibilityDocumentMap table
	if err := db.DB.AutoMigrate(&models.EligibilityDocumentMap{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	// Seed the database with default document types
	if err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.DefaultDocumentsRequired).Error; err != nil {
		log.Fatalf("Failed to seed database with default documents types: %s", err.Error())
	}

	// Start the HTTP server
	if err := r.Run(); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}
}
