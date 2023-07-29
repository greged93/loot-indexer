package indexer

import (
	"reflect"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
)

type IndexerConfig struct {
}

type Indexer struct {
	logger *log.Logger
}

func NewIndexerProducer() actor.Producer {
	return func() actor.Actor {
		return NewIndexer()
	}
}

func NewIndexer() *Indexer {
	return &Indexer{}
}

func (state *Indexer) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(ctx); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")
	case *actor.Stopping:
		if err := state.Clean(ctx); err != nil {
			state.logger.Error("error stopping actor", log.Error((err)))
			panic(err)
		}
		state.logger.Info("actor stopping")
	case *actor.Stopped:
		state.logger.Info("actor stopped")
	case *actor.Restarting:
		if err := state.Clean(ctx); err != nil {
			state.logger.Error("error restarting actor", log.Error((err)))
		}
	}
}

func (state *Indexer) Initialize(ctx actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", ctx.Self().Id),
		log.String("Type", reflect.TypeOf(*state).String()),
	)
	return nil
}

func (state *Indexer) Clean(ctx actor.Context) error {
	return nil
}
