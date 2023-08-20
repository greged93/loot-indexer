package loot_test

import (
	"loot-indexer/loot"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsPacking(t *testing.T) {
	// Given
	s := loot.Stats{
		Strength:     1,
		Dexterity:    2,
		Vitality:     3,
		Intelligence: 4,
		Wisdom:       5,
		Charisma:     6,
	}

	// When
	packed := loot.Pack(s)

	// Then
	assert.Equal(t, uint64(0x60504030201), packed)
}

func TestStatsUnPacking(t *testing.T) {
	//Given
	var s loot.Stats

	// When
	err := loot.Unpack(0x60504030201, &s)
	assert.NoError(t, err)

	// Then
	expected := loot.Stats{
		Strength:     1,
		Dexterity:    2,
		Vitality:     3,
		Intelligence: 4,
		Wisdom:       5,
		Charisma:     6,
	}
	assert.Equal(t, expected, s)
}
