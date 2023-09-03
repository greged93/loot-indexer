package tests

import (
	"context"
	starkclient "loot-indexer/stark-client"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	if cl == nil {
		t.Fatalf("nil client")
	}
	filter := starkclient.NewEventsArg(845505, 845505, nil, make([][]*felt.Felt, 0))
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

func TestGetBlock(t *testing.T) {
	if cl == nil {
		t.Fatalf("nil client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	block, err := cl.GetBlock(ctx)
	if !assert.Nil(t, err, "failed to get block: %v", err) {
		t.Fatal()
	}
	if !assert.NotZero(t, block, "expected non zero block: %v", err) {
		t.Fatal()
	}
}
