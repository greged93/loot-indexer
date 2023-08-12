package tests

import (
	"context"
	"fmt"
	starkclient "loot-indexer/stark-client"
	"os"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var cl *starkclient.Client

const (
	URL_GOERLI = "https://starknet-goerli.infura.io/v3/"
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
		fmt.Printf("failed to load env: %v", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestGetEvents(t *testing.T) {
	filter := starkclient.EventsArg{
		EventFilter: starkclient.EventFilter{
			FromBlock: &starkclient.BlockID{
				Pending: false, Latest: false, Hash: nil, Number: 845505,
			},
			ToBlock: &starkclient.BlockID{
				Pending: false, Latest: false, Hash: nil, Number: 845505,
			},
			Address: nil,
			Keys:    make([][]*felt.Felt, 0),
		},
		ResultPageRequest: starkclient.ResultPageRequest{
			ContinuationToken: "", ChunkSize: 32,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	events, err := cl.GetEvents(ctx, &filter)
	if !assert.Nil(t, err, "failed to get events: %v", err) {
		t.Fatal()
	}
	if !assert.NotEmpty(t, events.Events, "expected non empty events: %v", err) {
		t.Fatal()
	}
}
