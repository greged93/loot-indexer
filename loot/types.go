package loot

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

func FeltToBigInt(f *felt.Felt) (*big.Int, bool) {
	return new(big.Int).SetString(f.String(), 0)
}

type SqlBigInt big.Int

func (b *SqlBigInt) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	result := new(big.Int).SetBytes(bytes)
	*b = SqlBigInt(*result)
	return nil
}

func (sqlB SqlBigInt) Value() (driver.Value, error) {
	b := big.Int(sqlB)
	return json.Marshal(b.Bytes())
}

type Adventurer struct {
	AdventurerID        SqlBigInt
	LastAction          uint16        // 9 bits
	Health              uint16        // 9 bits
	Xp                  uint16        // 13 bits
	Stats               Stats         // 30 bits
	Gold                uint16        // 9 bits
	Weapon              ItemPrimitive // 21 bits
	Chest               ItemPrimitive // 21 bits
	Head                ItemPrimitive // 21 bits
	Waist               ItemPrimitive // 21 bits
	Foot                ItemPrimitive // 21 bits
	Hand                ItemPrimitive // 21 bits
	Neck                ItemPrimitive // 21 bits
	Ring                ItemPrimitive // 21 bits
	BeastHealth         uint16        // 9 bits
	StatPointsAvailable uint8         // 3 bits
	Mutated             bool          // not packed
}

type AdventurerMetadata struct {
	Name      *felt.Felt
	HomeRealm uint16
	Class     uint8
	Entropy   *felt.Felt
}

type ItemPrimitive struct {
	Id       uint8  // 7 bits
	Xp       uint16 // 9 bits
	Metadata uint8  // 5 bits
}

type Stats struct {
	Strength     uint8 // 5 bits
	Dexterity    uint8 // 5 bits
	Vitality     uint8 // 5 bits
	Intelligence uint8 // 5 bits
	Wisdom       uint8 // 5 bits
	Charisma     uint8 // 5 bits
}

type AdventurerState struct {
	Owner        *felt.Felt
	AdventurerId *felt.Felt
	Adventurer
}

type AdventurerStateWithBag struct {
	AdventurerState
	Bag
}

type Bag struct {
	Item1   ItemPrimitive
	Item2   ItemPrimitive
	Item3   ItemPrimitive
	Item4   ItemPrimitive
	Item5   ItemPrimitive
	Item6   ItemPrimitive
	Item7   ItemPrimitive
	Item8   ItemPrimitive
	Item9   ItemPrimitive
	Item10  ItemPrimitive
	Item11  ItemPrimitive
	Mutated bool
}

type CombatSpec struct {
	Tier     Tier
	ItemType Type
	Level    uint16
	Specials uint32 // packed bits of SpecialPowers
}

func (s *CombatSpec) Pack() uint64 {
	return uint64(s.Specials)<<32 | uint64(s.Level)<<16 | uint64(s.ItemType)<<8 | uint64(s.Tier)
}

type SpecialPowers struct {
	Special1 uint8
	Special2 uint8
	Special3 uint8
}

type LootWithPrice struct {
	Loot
	Price uint16
}

func (l *LootWithPrice) Pack() uint64 {
	return uint64(l.Price)<<32 | uint64(l.Loot.Pack())
}

type Loot struct {
	Id       uint8
	Tier     Tier
	ItemType Type
	Slot     Slot
}

func (l *Loot) Pack() uint32 {
	return uint32(l.Slot)<<24 | uint32(l.ItemType)<<16 | uint32(l.Tier)<<8 | uint32(l.Id)
}

type ItemSpecials struct {
	Special1 uint8 // 4 bit
	Special2 uint8 // 7 bits
	Special3 uint8 // 5 bits
}

func (s *ItemSpecials) Pack() uint32 {
	return uint32(s.Special3)<<16 | uint32(s.Special2)<<8 | uint32(s.Special1)
}
