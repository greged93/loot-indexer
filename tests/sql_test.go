package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, db)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestPostgresConnexion(t *testing.T) {
	_, err := InitDB()
	if !assert.Nil(t, err, "failed to connect to db: %v", err) {
		t.Fatal()
	}
}

func TestPostgresAddTable(t *testing.T) {
	db, err := InitDB()
	if !assert.Nil(t, err, "failed to connect to db: %v", err) {
		t.Fatal()
	}
	type SimpleTable struct {
		ID   uint `gorm:"primary_key"`
		Name string
	}

	err = db.AutoMigrate(&SimpleTable{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Create(&SimpleTable{Name: "test"}).Error
	if err != nil {
		log.Fatal(err)
	}
}
