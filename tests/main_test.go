package tests

import (
	"fmt"
	"loot-indexer/indexer"
	starkclient "loot-indexer/stark-client"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
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

	contractAddress, err := felt.Zero.SetString("0x1234")
	if err != nil {
		fmt.Printf("failed to convert 0x1234 to felt: %v", err)
		os.Exit(1)
	}

	db, err := InitDB()
	if err != nil {
		fmt.Printf("failed to initialize db: %v", err)
		os.Exit(1)
	}

	indexerConfig = indexer.NewIndexerConfig(url, contractAddress, db)

	os.Exit(m.Run())
}
