package tests

import (
	"loot-indexer/indexer"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

func TestIndexer(t *testing.T) {
	if indexerConfig == nil {
		t.Fatalf("nil config")
	}
	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)

	indexerPid, err := ctx.SpawnNamed(actor.PropsFromProducer(indexer.NewIndexerProducer(*indexerConfig), actor.WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))), "indexer")
	if !assert.Nil(t, err, "failed to spawn indexer: %v", err) {
		t.Fatal()
	}

	// sleep a little to be sure actor is ok
	time.Sleep(5 * time.Second)
	if err := ctx.PoisonFuture(indexerPid).Wait(); !assert.Nil(t, err, "failed to stop indexer: %v", err) {
		t.Fatal()
	}
}
