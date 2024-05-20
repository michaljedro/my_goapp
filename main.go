package main

import (
	"log"
	"net/http"

	"myapp/handlers"
	"myapp/utils"

	"github.com/gorilla/mux"
)

func main() {
    // Initialize the database
    db, err := utils.NewDB()
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    // Initialize the Vault client
    vaultClient, err := utils.NewVaultClient()
    if err != nil {
        log.Fatalf("Failed to connect to Vault: %v", err)
    }

    // Initialize the router
    r := mux.NewRouter()

    // Pass the database and Vault client to handlers
    handlers.InitializeHandlers(r, db, vaultClient)

    // Start the server
    log.Println("Server running on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
