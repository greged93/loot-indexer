package indexer

import (
	"fmt"
	"loot-indexer/loot"
	"reflect"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gorm.io/gorm"
)

type ErrorTracker struct {
	err error
}

type IndexerConfig struct {
	rpcUrl          string
	contractAddress *felt.Felt
	startBlock      uint64
	lastBlock       uint64
	db              *gorm.DB
}

func NewIndexerConfig(rpcUrl string, contractAddress *felt.Felt, start uint64, db *gorm.DB) *IndexerConfig {
	return &IndexerConfig{rpcUrl: rpcUrl, contractAddress: contractAddress, startBlock: start, db: db}
}

type Indexer struct {
	indexerConfig *IndexerConfig
	executor      *actor.PID
	logger        *log.Logger
}

func NewIndexerProducer(config *IndexerConfig) actor.Producer {
	return func() actor.Actor {
		return NewIndexer(config)
	}
}

func NewIndexer(config *IndexerConfig) *Indexer {
	return &Indexer{indexerConfig: config}
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
	case *ErrorTracker:
		state.logger.Error("error in child actor", log.Error(ctx.Message().(*ErrorTracker).err))
	}
}

func (state *Indexer) Initialize(ctx actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", ctx.Self().Id),
		log.String("Type", reflect.TypeOf(*state).String()),
	)

	err := state.indexerConfig.db.AutoMigrate(&loot.RawEvent{})
	if err != nil {
		return fmt.Errorf("error creating loot raw events table: %v", err)
	}

	err = state.indexerConfig.db.AutoMigrate(&loot.AdventurerStateWithBag{})
	if err != nil {
		return fmt.Errorf("error creating adventurer table: %v", err)
	}

	var event loot.RawEvent
	if err = state.indexerConfig.db.Order("block_number desc").First(&event).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error retrieving last event from db: %v", err)
	}
	if err == nil {
		if event.BlockNumber > state.indexerConfig.startBlock {
			state.indexerConfig.lastBlock = event.BlockNumber
		} else {
			state.indexerConfig.lastBlock = state.indexerConfig.startBlock
		}
	}

	props := actor.PropsFromProducer(NewExecutorProducer(state.indexerConfig.rpcUrl, state.indexerConfig.contractAddress, state.indexerConfig.lastBlock))
	pid, err := ctx.SpawnNamed(props, "executor")
	if err != nil {
		return fmt.Errorf("error spawning executor: %v", err)
	}
	state.executor = pid

	return nil
}

func (state *Indexer) Clean(ctx actor.Context) error {
	return nil
}
