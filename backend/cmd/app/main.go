package main

import (
	"github.com/drewbuiltit/trading-journal/backend/internal/auth"
	"github.com/drewbuiltit/trading-journal/backend/internal/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	s := store.NewMemoryStore()

	authHandler := &auth.AuthHandler{Store: s}

	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(auth.AuthMiddleWare)
	protected.HandleFunc("/", authHandler.ProtectedEndpoint).Methods("GET")

	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
