package main

import (
	"fmt"
	"github.com/drewbuiltit/trading-journal/backend/internal/auth"
	"github.com/drewbuiltit/trading-journal/backend/internal/models"
	"github.com/drewbuiltit/trading-journal/backend/internal/store"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect ot the database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Trade{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	auth.Init()

	router := mux.NewRouter()

	s := store.NewPostgresStore(db)

	authHandler := &auth.AuthHandler{Store: s}

	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(auth.AuthMiddleWare)
	protected.HandleFunc("/", authHandler.ProtectedEndpoint).Methods("GET")

	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
