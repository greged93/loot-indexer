package test

import (
	"context"
	starkclient "loot-indexer/stark-client"
	"os"
	"testing"
	"time"

	junoRpc "github.com/NethermindEth/juno/rpc"

	"github.com/joho/godotenv"
)

var cl *starkclient.Client

const (
	URL_GOERLI = "https://starknet-goerli.infura.io/v3/"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	infuraKey := os.Getenv("INFURA_KEY")
	url := URL_GOERLI + infuraKey
	cl, err = starkclient.Dial(url)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestGetEvents(t *testing.T) {
	filter := junoRpc.EventsArg{}
	ctx, err := context.WithTimeout(context.Background(), 15*time.Second)
	if err != nil {
		t.Errorf("error initializing context: %v", err)
	}
	cl.GetEvents(ctx, filter)
}
