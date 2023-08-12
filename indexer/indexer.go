package indexer

import (
	"fmt"
	"loot-indexer/loot"
	starkclient "loot-indexer/stark-client"
	"reflect"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gorm.io/gorm"
)

type IndexerConfig struct {
	rpcUrl          string
	contractAddress *felt.Felt
	lastBlock       uint64
	db              *gorm.DB
}

func NewIndexerConfig(rpcUrl string, contractAddress *felt.Felt, db *gorm.DB) *IndexerConfig {
	return &IndexerConfig{rpcUrl: rpcUrl, contractAddress: contractAddress, db: db}
}

type Indexer struct {
	indexerConfig IndexerConfig
	client        *starkclient.Client
	logger        *log.Logger
}

func NewIndexerProducer(config IndexerConfig) actor.Producer {
	return func() actor.Actor {
		return NewIndexer(config)
	}
}

func NewIndexer(config IndexerConfig) *Indexer {
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
	}
}

func (state *Indexer) Initialize(ctx actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", ctx.Self().Id),
		log.String("Type", reflect.TypeOf(*state).String()),
	)

	err := state.indexerConfig.db.AutoMigrate(&loot.LootEvent{})
	if err != nil {
		return fmt.Errorf("error creating loot table: %v", err)
	}

	var event loot.LootEvent
	if err = state.indexerConfig.db.Order("block_number desc").First(&event).Error; err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error retrieving last event from db: %v", err)
	}
	if err == nil {
		state.indexerConfig.lastBlock = event.BlockNumber
	}

	client, err := starkclient.Dial(state.indexerConfig.rpcUrl)
	if err != nil {
		return fmt.Errorf("error dialing starknet client at url %s: %v", state.indexerConfig.rpcUrl, err)
	}
	state.client = client

	return nil
}

func (state *Indexer) Clean(ctx actor.Context) error {
	return nil
}
