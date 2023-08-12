package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/joho/godotenv"
)

const (
	URL_GOERLI = "https://starknet-goerli.infura.io/v3/"
)

var (
	indexer *actor.PID
	done    = make(chan os.Signal)
)

type GuardActor struct{}

func (guard *GuardActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		ctx.Watch(indexer)
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
	infuraKey := os.Getenv("URL_GOERLI")
	url := URL_GOERLI + infuraKey
	fmt.Printf("goerli url %s", url)
}
