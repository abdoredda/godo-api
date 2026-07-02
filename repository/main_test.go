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

	db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to open test db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping test db: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("failed to truncate tasks table: %v", err)
	}

	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}
