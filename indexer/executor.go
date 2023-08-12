package indexer

import (
	"fmt"
	starkclient "loot-indexer/stark-client"
	"reflect"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
)

type Executor struct {
	rpcUrl          string
	contractAddress *felt.Felt
	client          *starkclient.Client
	logger          *log.Logger
}

func NewExecutorProducer(url string, contractAddress *felt.Felt) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(url, contractAddress)
	}
}

func NewExecutor(url string, contractAddress *felt.Felt) actor.Actor {
	return &Executor{rpcUrl: url, contractAddress: contractAddress}
}

func (state *Executor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(ctx); err != nil {
			state.logger.Error("error initializing actor", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")
	case *actor.Stopping:
		if err := state.Clean(ctx); err != nil {
			state.logger.Error("error stopping actor", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")
	case *actor.Stopped:
		state.logger.Info("actor stopped")
	case *actor.Restarting:
		if err := state.Clean(ctx); err != nil {
			state.logger.Error("error restarting actor", log.Error(err))
		}
		state.logger.Info("actor restarting")
	}
}

func (state *Executor) Initialize(ctx actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", ctx.Self().Id),
		log.String("Type", reflect.TypeOf(*state).String()),
	)

	client, err := starkclient.Dial(state.rpcUrl)
	if err != nil {
		return fmt.Errorf("error dialing starknet client at url %s: %v", state.rpcUrl, err)
	}
	state.client = client

	return nil
}

func (state *Executor) Clean(ctx actor.Context) error {
	return nil
}
