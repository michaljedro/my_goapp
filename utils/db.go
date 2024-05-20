package utils

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/ClickHouse/clickhouse-go"
)

func NewDB() (*sqlx.DB, error) {
    db, err := sqlx.Open("clickhouse", "tcp://clickhouse:9000?debug=true")
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    log.Println("Connected to ClickHouse database")
    return db, nil
}
