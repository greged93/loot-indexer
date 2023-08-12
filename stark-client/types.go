package starkclient

import (
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
)

type EventsArg struct {
	EventFilter
	ResultPageRequest
}

type EventFilter struct {
	FromBlock *BlockID       `json:"from_block"`
	ToBlock   *BlockID       `json:"to_block"`
	Address   *felt.Felt     `json:"address"`
	Keys      [][]*felt.Felt `json:"keys"`
}

type BlockID struct {
	Pending bool
	Latest  bool
	Hash    *felt.Felt
	Number  uint64
}

type ResultPageRequest struct {
	ContinuationToken string `json:"continuation_token,omitempty"`
	ChunkSize         uint64 `json:"chunk_size" validate:"min=1"`
}

func (blockId *BlockID) MarshalJSON() ([]byte, error) {
	type BlockTag struct {
		BlockTag string `json:"block_tag"`
	}
	type BlockHash struct {
		BlockHash *felt.Felt `json:"block_hash"`
	}
	type BlockNumber struct {
		BlockNumber uint64 `json:"block_number"`
	}
	if blockId.Pending {
		return json.Marshal(BlockTag{BlockTag: "pending"})
	}
	if blockId.Latest {
		return json.Marshal(BlockTag{BlockTag: "latest"})
	}
	if blockId.Hash != nil {
		return json.Marshal(BlockHash{BlockHash: blockId.Hash})
	}
	return json.Marshal(BlockNumber{BlockNumber: blockId.Number})
}
