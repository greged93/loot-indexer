package indexer

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
)

type IndexerConfig struct {
	rpcUrl    string
	sqlConfig SqlConfig
}

type SqlConfig struct {
	Host     string
	Port     uint
	Db       string
	User     string
	Password string
}

type Indexer struct {
	indexerConfig IndexerConfig
	db            *sql.DB
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

	sqlConfig := state.indexerConfig.sqlConfig
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		sqlConfig.Host, sqlConfig.Port, sqlConfig.User, sqlConfig.Password, sqlConfig.Db)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("error opening sql database: %v", err)
	}
	state.db = db
	return nil
}

func (state *Indexer) Clean(ctx actor.Context) error {
	return nil
}
