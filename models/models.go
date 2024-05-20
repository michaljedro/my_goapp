package models

type User struct {
    Username string `json:"username" db:"username"`
    Email    string `json:"email" db:"email"`
}

type Key struct {
    Name     string `json:"name" db:"name"`
    Secret   string `json:"secret" db:"-"`
    Username string `json:"username" db:"username"`
}

type Item struct {
    Sector      string `json:"sector" db:"sector"`
    Name        string `json:"name" db:"name"`
    Description string `json:"description" db:"description"`
}

type Policy struct {
    Name       string `json:"name" db:"name"`
    Definition string `json:"definition" db:"definition"`
}
