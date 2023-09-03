package starkclient

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	junoRpc "github.com/NethermindEth/juno/rpc"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	Client *rpc.Client
}

func Dial(url string) (*Client, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (*Client, error) {
	client, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return NewClient(client), nil
}

func NewClient(client *rpc.Client) *Client {
	return &Client{Client: client}
}

func (client *Client) GetEvents(ctx context.Context, filter *EventsArg) (*junoRpc.EventsChunk, error) {
	res := &junoRpc.EventsChunk{}
	err := client.Client.CallContext(ctx, res, "starknet_getEvents", filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *Client) GetBlock(ctx context.Context) (uint64, error) {
	res := &felt.Felt{}
	err := client.Client.CallContext(ctx, res, "starknet_blockNumber")
	if err != nil {
		return 0, err
	}
	return res.Impl().Uint64(), nil
}
