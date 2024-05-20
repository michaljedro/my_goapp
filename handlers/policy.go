package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func initializePolicyHandlers(r *mux.Router, db *sqlx.DB) {
    r.HandleFunc("/", createPolicy(db)).Methods("POST")
    r.HandleFunc("/{name}", getPolicy(db)).Methods("GET")
    r.HandleFunc("/{name}", deletePolicy(db)).Methods("DELETE")
    r.HandleFunc("/", listPolicies(db)).Methods("GET")
}

func createPolicy(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var policy models.Policy
        if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err := db.NamedExec(`INSERT INTO policies (name, definition) VALUES (:name, :definition)`, &policy)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func getPolicy(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        name := vars["name"]

        var policy models.Policy
        err := db.Get(&policy, "SELECT name, definition FROM policies WHERE name=$1", name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(policy)
    }
}

func deletePolicy(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        name := vars["name"]

        _, err := db.Exec("DELETE FROM policies WHERE name=$1", name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func listPolicies(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var policies []models.Policy
        err := db.Select(&policies, "SELECT name, definition FROM policies")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(policies)
    }
}
