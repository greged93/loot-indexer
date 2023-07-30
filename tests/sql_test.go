package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	db       = "loot"
	user     = "loot"
	password = "loot-survivor"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func BeforeEach(t *testing.T) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, db)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestPostgresConnexion(t *testing.T) {
	db := BeforeEach(t)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostgresAddTable(t *testing.T) {
	db := BeforeEach(t)
	defer db.Close()

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS loot (id SERIAL PRIMARY KEY, name VARCHAR(50))")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO loot (name) VALUES ($1)", "hero")
	if err != nil {
		log.Fatal(err)
	}
}
