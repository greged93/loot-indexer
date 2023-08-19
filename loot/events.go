package loot

import (
	"github.com/NethermindEth/juno/core/felt"
)

type EventType uint8

const (
	StarGameType EventType = iota
	DiscoveredHealthType
	DiscoveredGoldType
	DiscoveredXPType
	DodgedObstacleType
	HitByObstacleType
	DiscoveredBeastType
	AmbushedByBeastType
	AttackedBeastType
	AttackedByBeastType
	SlayedBeastType
	FleeFailedType
	FleeSucceededType
	PurchasedPotionsType
	PurchasedItemsType
	EquippedItemsType
	DroppedItemsType
	ItemLeveledUpType
	ItemSpecialUnlockedType
	NewHighScoreType
	AdventurerDiedType
	ShopAvailableType
	AdventurerLeveledUpType
	NewItemsAvailableType
	IdleDamagePenaltyType
	UpgradeAvailableType
	AdventurerUpgradedType
)

type RawEvent struct {
	Type        EventType `json:"type"`
	Event       []byte    `json:"event"`
	BlockNumber uint64    `json:"block"`
}

type StartGame struct {
	AdventurerState
	AdventurerMetadata
}

type DiscoveredHealth struct {
	AdventurerState
	HealthAmount uint16
}

type DiscoveredGold struct {
	AdventurerState
	Goldamount uint16
}

type DiscoveredXP struct {
	AdventurerState
	XpAmount uint16
}

type DodgedObstacle struct {
	AdventurerState
	Id                 uint8
	Level              uint16
	DamageTaken        uint16
	DamageLocation     uint8
	XpEarnedAdventurer uint16
	XpEarnedItems      uint16
}

type HitByObstacle struct {
	AdventurerState    AdventurerState
	Id                 uint8
	Level              uint16
	DamageTaken        uint16
	DamageLocation     uint8
	XpEarnedAdventurer uint16
	XpEarnedItems      uint16
}

type DiscoveredBeast struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

type AmbushedByBeast struct {
	AdventurerState
	Seed        *felt.Felt
	Id          uint8
	BeastSpecs  CombatSpec
	Damage      uint16
	CriticalHit bool
	Location    uint8
}

type AttackedBeast struct {
	AdventurerState
	Seed        *felt.Felt
	Id          uint8
	BeastSpecs  CombatSpec
	Damage      uint16
	CriticalHit bool
	Location    uint8
}

type AttackedByBeast struct {
	AdventurerState
	Seed        *felt.Felt
	Id          uint8
	BeastSpecs  CombatSpec
	Damage      uint16
	CriticalHit bool
	Location    uint8
}

type SlayedBeast struct {
	AdventurerState
	Seed               *felt.Felt
	Id                 uint8
	BeastSpecs         CombatSpec
	DamageDealt        uint16
	CriticalHit        bool
	XpEarnedAdventurer uint16
	XpEarnedItems      uint16
	GoldEarned         uint16
}

type FleeFailed struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

type FleeSucceeded struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

type PurchasedPotions struct {
	AdventurerState
	Quantity uint8
	Cost     uint16
	Health   uint16
}

type PurchasedItems struct {
	AdventurerStateWithBag
	Purchases []*LootWithPrice
}

type EquippedItems struct {
	AdventurerStateWithBag
	EquippedItems   []uint8
	UnequippedItems []uint8
}

type DroppedItems struct {
	AdventurerStateWithBag
	ItemIds []uint8
}

type ItemLeveledUp struct {
	AdventurerState
	ItemId        uint8
	PreviousLevel uint8
	NewLevel      uint8
}

type ItemSpecialUnlocked struct {
	AdventurerState
	Id       uint8
	Level    uint8
	Specials ItemSpecials
}

type NewHighScore struct {
	AdventurerState
	Rank uint8 // 1-3
}

type AdventurerDied struct {
	AdventurerState
	KilledByBeast    bool
	KilledByoBstacle bool
	KillerId         uint8
}

type ShopAvailable struct {
	Inventory []*Loot
}

type AdventurerLeveledUp struct {
	AdventurerState
	PreviousLevel uint8
	NewLevel      uint8
}

type NewItemsAvailable struct {
	AdventurerState
	Items []*LootWithPrice
}

type IdleDamagePenalty struct {
	AdventurerState
	IdleBlocks  uint16
	DamageTaken uint16
}

type UpgradeAvailable struct {
	AdventurerState
}

type AdventurerUpgraded struct {
	AdventurerStateWithBag
	StrengthIncrease     uint8
	DexterityIncrease    uint8
	VitalityIncrease     uint8
	IntelligenceIncrease uint8
	WisdomIncrease       uint8
	CharismaIncrease     uint8
}
