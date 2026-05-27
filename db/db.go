package db

import "database/sql"

// sql.Open doesn't actually connect — it just validates the config and returns a handle. The real connection happens on the first query
func NewDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
