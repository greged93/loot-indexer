package tests

import (
	"fmt"
	"loot-indexer/indexer"
	starkclient "loot-indexer/stark-client"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const (
	host     = "localhost"
	port     = 5432
	db       = "loot"
	user     = "loot"
	password = "loot-survivor"

	URL_GOERLI = "https://starknet-goerli.infura.io/v3/"
)

var (
	cl            *starkclient.Client
	indexerConfig *indexer.IndexerConfig
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("failed to load env: %v", err)
		os.Exit(1)
	}
	infuraKey := os.Getenv("INFURA_KEY")
	url := URL_GOERLI + infuraKey

	cl, err = starkclient.Dial(url)
	if err != nil {
		fmt.Printf("failed to dial client env: %v", err)
		os.Exit(1)
	}

	indexerConfig = indexer.NewIndexerConfig(url, indexer.SqlConfig{Host: host, Port: port, Db: db, User: user, Password: password})

	os.Exit(m.Run())
}
