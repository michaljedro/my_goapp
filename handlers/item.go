package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"

	"github.com/gorilla/mux"
	"github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
)

func initializeItemHandlers(r *mux.Router, db *sqlx.DB, vaultClient *api.Client) {
    r.HandleFunc("/", createItem(db, vaultClient)).Methods("POST")
    r.HandleFunc("/{sector}/{name}", getItem(db)).Methods("GET")
    r.HandleFunc("/{sector}/{name}", deleteItem(db)).Methods("DELETE")
    r.HandleFunc("/", listItems(db)).Methods("GET")
}

func createItem(db *sqlx.DB, vaultClient *api.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var item models.Item
        if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err := db.NamedExec(`INSERT INTO items (sector, name, description) VALUES (:sector, :name, :description)`, &item)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func getItem(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        sector := vars["sector"]
        name := vars["name"]

        var item models.Item
        err := db.Get(&item, "SELECT sector, name, description FROM items WHERE sector=$1 AND name=$2", sector, name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(item)
    }
}

func deleteItem(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        sector := vars["sector"]
        name := vars["name"]

        _, err := db.Exec("DELETE FROM items WHERE sector=$1 AND name=$2", sector, name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func listItems(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var items []models.Item
        err := db.Select(&items, "SELECT sector, name, description FROM items")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(items)
    }
}
