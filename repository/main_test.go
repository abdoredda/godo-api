package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to open test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping test db: %v", err)
	}
	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}
