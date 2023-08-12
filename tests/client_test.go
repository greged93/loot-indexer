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
