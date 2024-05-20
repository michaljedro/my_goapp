package handlers

import (
    "encoding/json"
    "net/http"

    "myapp/models"

    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
)

func initializeUserHandlers(r *mux.Router, db *sqlx.DB) {
    r.HandleFunc("/", createUser(db)).Methods("POST")
    r.HandleFunc("/{username}", getUser(db)).Methods("GET")
    r.HandleFunc("/{username}", deleteUser(db)).Methods("DELETE")
    r.HandleFunc("/", listUsers(db)).Methods("GET")
}

func createUser(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err := db.NamedExec(`INSERT INTO users (username, email) VALUES (:username, :email)`, &user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func getUser(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        username := vars["username"]

        var user models.User
        err := db.Get(&user, "SELECT username, email FROM users WHERE username=$1", username)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(user)
    }
}

func deleteUser(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        username := vars["username"]

        _, err := db.Exec("DELETE FROM users WHERE username=$1", username)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func listUsers(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var users []models.User
        err := db.Select(&users, "SELECT username, email FROM users")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(users)
    }
}
