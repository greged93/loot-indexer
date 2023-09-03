package indexer

import (
	goContext "context"
	"fmt"
	starkclient "loot-indexer/stark-client"
	"reflect"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	junoRpc "github.com/NethermindEth/juno/rpc"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
)

type Executor struct {
	rpcUrl          string
	contractAddress *felt.Felt
	startBlock      uint64
	client          *starkclient.Client
	logger          *log.Logger
}

func NewExecutorProducer(url string, contractAddress *felt.Felt, start uint64) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(url, contractAddress, start)
	}
}

func NewExecutor(url string, contractAddress *felt.Felt, start uint64) actor.Actor {
	return &Executor{rpcUrl: url, contractAddress: contractAddress, startBlock: start}
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
	case *junoRpc.EventsChunk:
		// TODO handle events
		state.logger.Info("received events chunk", log.Object("events", ctx.Message().(*junoRpc.EventsChunk).Events))
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

	// Launch the indexer
	go func(parent *actor.PID, self *actor.PID) {
		var lastBlock = state.startBlock
		for {
			context, cancel := goContext.WithTimeout(goContext.Background(), 5*time.Second)
			currentBlock, err := state.client.GetBlock(context)
			if err != nil {
				state.logger.Error("error getting current block", log.Error(err))
				cancel()
				return
			}

			if currentBlock > lastBlock {
				filter := starkclient.NewEventsArg(lastBlock, lastBlock+1000, state.contractAddress, make([][]*felt.Felt, 0))
				for {
					events, err := client.GetEvents(context, &filter)
					if err != nil {
						ctx.Send(parent, &ErrorTracker{err: err})
						cancel()
						return
					}
					ctx.Send(self, events)
					if events.ContinuationToken == "" {
						break
					}
					filter.ContinuationToken = events.ContinuationToken
				}
				lastBlock += 1000
				cancel()
			} else {
				cancel()
				return
			}
		}
	}(ctx.Parent(), ctx.Self())

	return nil
}

func (state *Executor) Clean(ctx actor.Context) error {
	return nil
}
