package loot

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/jackc/pgx/v5/pgtype"
)

type LootEvent struct {
	// Owner
	Owner        pgtype.Numeric `gorm:"primaryKey,type:numeric"`
	AdventurerId pgtype.Numeric `gorm:"primaryKey,type:numeric"`
	// Block
	BlockNumber uint64
	// Event
	Type Event
	Seed pgtype.Numeric
	Id   uint8
	// Earned
	HealthAmount       uint16
	GoldAmount         uint16
	GoldEarned         uint16
	XpAmount           uint16
	Level              uint16
	XpEarnedAdventurer uint16
	XpEarnedItems      uint16
	// Combat
	BeastSpecs     uint64
	Location       uint8
	Damage         uint16
	CriticalHit    bool
	IdleBlocks     uint16
	DamageTaken    uint16
	DamageDealt    uint16
	DamageLocation uint8
	// Purchased
	Quantity  uint8
	Cost      uint16
	Health    uint16
	Purchases []uint64 `gorm:"type:bigint[]"`
	// Equipped
	EquippedItems   []uint8 `gorm:"type:bigint[]"`
	UnequippedItems []uint8 `gorm:"type:bigint[]"`
	ItemIds         []uint8 `gorm:"type:bigint[]"`
	// Level up
	ItemId        uint8
	PreviousLevel uint8
	NewLevel      uint8
	Specials      uint32
	Rank          uint8 // 1-3
	// Killed
	KilledByBeast    bool
	KilledByoBstacle bool
	KillerId         uint8
	// Shop
	Inventory []uint32 `gorm:"type:bigint[]"`
	Items     []uint64 `gorm:"type:bigint[]"`
	// Upgrade
	StrengthIncrease     uint8
	DexterityIncrease    uint8
	VitalityIncrease     uint8
	IntelligenceIncrease uint8
	WisdomIncrease       uint8
	CharismaIncrease     uint8
}

type StartGame struct {
	AdventurerState
	AdventurerMetadata
}

func (e *LootEvent) FromDiscoveredHealthEvent(event *DiscoveredHealth) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DiscoveredHealthEvent
	e.HealthAmount = event.HealthAmount
	return nil
}

type DiscoveredHealth struct {
	AdventurerState
	HealthAmount uint16
}

func (e *LootEvent) FromDiscoveredGoldEvent(event *DiscoveredGold) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DiscoveredGoldEvent
	e.GoldAmount = event.Goldamount
	return nil
}

type DiscoveredGold struct {
	AdventurerState
	Goldamount uint16
}

func (e *LootEvent) FromDiscoveredXPEvent(event *DiscoveredXP) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DiscoveredXPEvent
	e.XpAmount = event.XpAmount
	return nil
}

type DiscoveredXP struct {
	AdventurerState
	XpAmount uint16
}

func (e *LootEvent) FromDodgedObstacleEvent(event *DodgedObstacle) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DodgedObstacleEvent
	e.Id = event.Id
	e.Level = event.Level
	e.DamageTaken = event.DamageTaken
	e.DamageLocation = event.DamageLocation
	e.XpEarnedAdventurer = event.XpEarnedAdventurer
	e.XpEarnedItems = event.XpEarnedItems
	return nil
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

func (e *LootEvent) FromHitByObstacleEvent(event *HitByObstacle) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = HitByObstacleEvent
	e.Id = event.Id
	e.Level = event.Level
	e.DamageTaken = event.DamageTaken
	e.DamageLocation = event.DamageLocation
	e.XpEarnedAdventurer = event.XpEarnedAdventurer
	e.XpEarnedItems = event.XpEarnedItems
	return nil
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

func (e *LootEvent) FromDiscoveredBeastEvent(event *DiscoveredBeast) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DiscoveredBeastEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	return nil
}

type DiscoveredBeast struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

func (e *LootEvent) FromAmbushedByBeastEvent(event *AmbushedByBeast) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AmbushedByBeastEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	e.Damage = event.Damage
	e.CriticalHit = event.CriticalHit
	e.Location = event.Location
	return nil
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

func (e *LootEvent) FromAttackedBeastEvent(event *AttackedBeast) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AttackedBeastEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	e.Damage = event.Damage
	e.CriticalHit = event.CriticalHit
	e.Location = event.Location
	return nil
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

func (e *LootEvent) FromAttackedByBeastEvent(event *AttackedByBeast) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AttackedByBeastEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	e.Damage = event.Damage
	e.CriticalHit = event.CriticalHit
	e.Location = event.Location
	return nil
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

func (e *LootEvent) FromSlayedBeastEvent(event *SlayedBeast) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = SlayedBeastEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	e.DamageDealt = event.DamageDealt
	e.CriticalHit = event.CriticalHit
	e.XpEarnedAdventurer = event.XpEarnedAdventurer
	e.XpEarnedItems = event.XpEarnedItems
	e.GoldEarned = event.GoldEarned
	return nil
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

func (e *LootEvent) FromFleeFailedEvent(event *FleeFailed) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = FleeFailedEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	return nil
}

type FleeFailed struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

func (e *LootEvent) FromFleeSucceededEvent(event *FleeSucceeded) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = FleeSucceededEvent
	e.Id = event.Id
	e.BeastSpecs = event.BeastSpecs.Pack()
	return nil
}

type FleeSucceeded struct {
	AdventurerState
	Seed       *felt.Felt
	Id         uint8
	BeastSpecs CombatSpec
}

func (e *LootEvent) FromPurchasedPotionsEvent(event *PurchasedPotions) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = PurchasedPotionsEvent
	e.Quantity = event.Quantity
	e.Cost = event.Cost
	e.Health = event.Health
	return nil
}

type PurchasedPotions struct {
	AdventurerState
	Quantity uint8
	Cost     uint16
	Health   uint16
}

func (e *LootEvent) FromPurchasedItemsEvent(event *PurchasedItems) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = PurchasedItemsEvent
	var purchases []uint64
	for _, p := range event.Purchases {
		purchases = append(purchases, p.Pack())
	}
	e.Purchases = purchases
	return nil
}

type PurchasedItems struct {
	AdventurerStateWithBag
	Purchases []*LootWithPrice
}

func (e *LootEvent) FromEquippedItemsEvent(event *EquippedItems) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = EquippedItemsEvent
	e.EquippedItems = event.EquippedItems
	e.UnequippedItems = event.UnequippedItems
	return nil
}

type EquippedItems struct {
	AdventurerStateWithBag
	EquippedItems   []uint8
	UnequippedItems []uint8
}

func (e *LootEvent) FromDroppedItemsEvent(event *DroppedItems) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = DroppedItemsEvent
	e.ItemIds = event.ItemIds
	return nil
}

type DroppedItems struct {
	AdventurerStateWithBag
	ItemIds []uint8
}

func (e *LootEvent) FromItemLeveledUpEvent(event *ItemLeveledUp) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = ItemLeveledUpEvent
	e.ItemId = event.ItemId
	e.PreviousLevel = event.PreviousLevel
	e.NewLevel = event.NewLevel
	return nil
}

type ItemLeveledUp struct {
	AdventurerState
	ItemId        uint8
	PreviousLevel uint8
	NewLevel      uint8
}

func (e *LootEvent) FromItemSpecialUnlockedEvent(event *ItemSpecialUnlocked) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = ItemSpecialUnlockedEvent
	e.Id = event.Id
	e.Level = uint16(event.Level)

	e.Specials = event.Specials.Pack()
	return nil
}

type ItemSpecialUnlocked struct {
	AdventurerState
	Id       uint8
	Level    uint8
	Specials ItemSpecials
}

func (e *LootEvent) FromNewHighScoreEvent(event *NewHighScore) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = NewHighScoreEvent
	e.Rank = event.Rank
	return nil
}

type NewHighScore struct {
	AdventurerState
	Rank uint8 // 1-3
}

func (e *LootEvent) FromAdventurerDiedEvent(event *AdventurerDied) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AdventurerDiedEvent
	e.KilledByBeast = event.KilledByBeast
	e.KilledByoBstacle = event.KilledByoBstacle
	e.KillerId = event.KillerId
	return nil
}

type AdventurerDied struct {
	AdventurerState
	KilledByBeast    bool
	KilledByoBstacle bool
	KillerId         uint8
}

func (e *LootEvent) FromShopAvailableEvent(event *ShopAvailable) {
	e.Type = ShopAvailableEvent

	var inventory []uint32
	for _, l := range event.Inventory {
		inventory = append(inventory, l.Pack())
	}
	e.Inventory = inventory
}

type ShopAvailable struct {
	Inventory []*Loot
}

func (e *LootEvent) FromAdventurerLeveledUpEvent(event *AdventurerLeveledUp) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AdventurerLeveledUpEvent
	e.PreviousLevel = event.PreviousLevel
	e.NewLevel = event.NewLevel
	return nil
}

type AdventurerLeveledUp struct {
	AdventurerState
	PreviousLevel uint8
	NewLevel      uint8
}

func (e *LootEvent) FromNewItemsAvailableEvent(event *NewItemsAvailable) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = NewItemsAvailableEvent

	var items []uint64
	for _, l := range event.Items {
		items = append(items, l.Pack())
	}
	e.Items = items
	return nil
}

type NewItemsAvailable struct {
	AdventurerState
	Items []*LootWithPrice
}

func (e *LootEvent) FromIdleDamagePenaltyEvent(event *IdleDamagePenalty) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = IdleDamagePenaltyEvent
	e.IdleBlocks = event.IdleBlocks
	e.DamageTaken = event.DamageTaken
	return nil
}

type IdleDamagePenalty struct {
	AdventurerState
	IdleBlocks  uint16
	DamageTaken uint16
}

func (e *LootEvent) FromUpgradeAvailableEvent(event *UpgradeAvailable) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = UpgradeAvailableEvent
	return nil
}

type UpgradeAvailable struct {
	AdventurerState
}

func (e *LootEvent) FromAdventurerUpgradedEvent(event *AdventurerUpgraded) error {
	var err error
	e.Owner, err = FeltToNumeric(event.AdventurerState.Owner)
	if err != nil {
		return fmt.Errorf("failed to convert owner to numeric: %w", err)
	}

	e.Type = AdventurerUpgradedEvent
	e.StrengthIncrease = event.StrengthIncrease
	e.DexterityIncrease = event.DexterityIncrease
	e.VitalityIncrease = event.VitalityIncrease
	e.IntelligenceIncrease = event.IntelligenceIncrease
	e.WisdomIncrease = event.WisdomIncrease
	e.CharismaIncrease = event.CharismaIncrease
	return nil
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
