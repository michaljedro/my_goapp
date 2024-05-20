package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"

	"github.com/gorilla/mux"
	"github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
)

func initializeKeyHandlers(r *mux.Router, db *sqlx.DB, vaultClient *api.Client) {
    r.HandleFunc("/", createKey(db, vaultClient)).Methods("POST")
    r.HandleFunc("/{username}/{name}", deleteKey(db)).Methods("DELETE")
    r.HandleFunc("/", listKeys(db)).Methods("GET")
}

func createKey(db *sqlx.DB, vaultClient *api.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var key models.Key
        if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Store secret in Vault
        secretPath := "secret/data/keys/" + key.Name
        _, err := vaultClient.Logical().Write(secretPath, map[string]interface{}{
            "data": map[string]interface{}{
                "secret": key.Secret,
            },
        })
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        _, err = db.NamedExec(`INSERT INTO keys (name, username) VALUES (:name, :username)`, &key)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func deleteKey(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        username := vars["username"]
        name := vars["name"]

        _, err := db.Exec("DELETE FROM keys WHERE username=$1 AND name=$2", username, name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func listKeys(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var keys []models.Key
        err := db.Select(&keys, "SELECT name, username FROM keys")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(keys)
    }
}
