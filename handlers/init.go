package handlers

import (
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    "github.com/hashicorp/vault/api"
)

var DB *sqlx.DB
var VaultClient *api.Client

func InitializeHandlers(r *mux.Router, db *sqlx.DB, vaultClient *api.Client) {
    DB = db
    VaultClient = vaultClient

    // Define your routes here
    r.HandleFunc("/users", CreateUser).Methods("POST")
    r.HandleFunc("/keys", CreateKey).Methods("POST")
    r.HandleFunc("/items", CreateItem).Methods("POST")
    r.HandleFunc("/policies", CreatePolicy).Methods("POST")
}
