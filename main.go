package main

import (
	"fmt"
	"loot-indexer/indexer"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	URL_GOERLI = "https://starknet-goerli.infura.io/v3/"
)

var (
	idx  *actor.PID
	done = make(chan os.Signal)
)

type GuardActor struct{}

func (guard *GuardActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		ctx.Watch(idx)
	case actor.Terminated:
		done <- syscall.SIGTERM
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("failed to load env: %v", err)
		os.Exit(1)
	}

	infuraKey := loadVariable("INFURA_KEY")
	url := URL_GOERLI + infuraKey

	addr := loadVariable("CONTRACT_ADDRESS")
	db := loadVariable("DB")
	user := loadVariable("USER")
	password := loadVariable("PASSWORD")
	startBlock, err := strconv.ParseUint(loadVariable("START_BLOCK"), 10, 32)
	if err != nil {
		panic(err)
	}

	address, err := felt.Zero.SetString(addr)
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable",
		user, password, db)
	gormDB, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cfg := indexer.NewIndexerConfig(url, address, startBlock, gormDB)

	as := actor.NewActorSystem()
	ctx := actor.NewRootContext(as, nil)

	idx, err := ctx.SpawnNamed(actor.PropsFromProducer(indexer.NewIndexerProducer(cfg), actor.WithSupervisor(actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))), "indexer")
	if err != nil {
		panic(err)
	}

	guard, err := ctx.SpawnNamed(actor.PropsFromProducer(func() actor.Actor { return &GuardActor{} }), "guard")
	if err != nil {
		panic(err)
	}
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	tries := 0
	_ = ctx.PoisonFuture(guard).Wait()
	for {
		if err := ctx.PoisonFuture(idx).Wait(); err != nil {
			break
		}
		tries++
		if tries > 80 {
			panic("failed to stop indexer")
		}
	}
}

func loadVariable(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("failed to load %s", name))
	}
	return value
}
