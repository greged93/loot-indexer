package starkclient_test

import (
	starkclient "loot-indexer/stark-client"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMarshallBlockIdPending(t *testing.T) {
	id := starkclient.BlockID{Pending: true}
	b, err := id.MarshalJSON()
	if !assert.Nil(t, err, "failed to marshall block id: %v", err) {
		t.Fatal()
	}
	if !assert.Equal(t, `{"block_tag":"pending"}`, string(b), "incorrect marshalled tag") {
		t.Fatal()
	}
}

func TestMarshallBlockIdLatest(t *testing.T) {
	id := starkclient.BlockID{Latest: true}
	b, err := id.MarshalJSON()
	if !assert.Nil(t, err, "failed to marshall block id: %v", err) {
		t.Fatal()
	}
	if !assert.Equal(t, `{"block_tag":"latest"}`, string(b), "incorrect marshalled tag") {
		t.Fatal()
	}
}

func TestMarshallBlockIdHash(t *testing.T) {
	num, err := felt.Zero.SetString("0x1234")
	if !assert.Nil(t, err, "failed to set string: %v", err) {
		t.Fatal()
	}

	id := starkclient.BlockID{Hash: num}
	b, err := id.MarshalJSON()
	if !assert.Nil(t, err, "failed to marshall block id: %v", err) {
		t.Fatal()
	}

	if !assert.Equal(t, `{"block_hash":"0x1234"}`, string(b), "incorrect marshalled tag") {
		t.Fatal()
	}
}

func TestMarshallBlockIdNumber(t *testing.T) {
	id := starkclient.BlockID{Number: 100}
	b, err := id.MarshalJSON()
	if !assert.Nil(t, err, "failed to marshall block id: %v", err) {
		t.Fatal()
	}

	if !assert.Equal(t, `{"block_number":100}`, string(b), "incorrect marshalled tag") {
		t.Fatal()
	}
}
