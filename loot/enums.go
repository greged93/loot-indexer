package loot

type WeaponEffectiveness uint8

type Event uint8

const (
	StartGameEvent Event = iota
	DiscoveredHealthEvent
	DiscoveredGoldEvent
	DiscoveredXPEvent
	DodgedObstacleEvent
	HitByObstacleEvent
	DiscoveredBeastEvent
	AmbushedByBeastEvent
	AttackedBeastEvent
	AttackedByBeastEvent
	SlayedBeastEvent
	FleeFailedEvent
	FleeSucceededEvent
	PurchasedPotionsEvent
	PurchasedItemsEvent
	EquippedItemsEvent
	DroppedItemsEvent
	ItemLeveledUpEvent
	ItemSpecialUnlockedEvent
	NewHighScoreEvent
	AdventurerDiedEvent
	ShopAvailableEvent
	AdventurerLeveledUpEvent
	NewItemsAvailableEvent
	IdleDamagePenaltyEvent
	UpgradeAvailableEvent
	AdventurerUpgradedEvent
)

const (
	Weak WeaponEffectiveness = iota
	Fair
	Strong
)

type Type uint8

const (
	NoneType Type = iota
	MagicOrClothType
	BladeOrHideType
	BludgeonOrMetalType
	NecklaceType
	RingType
)

type Tier uint8

const (
	NoneTier Tier = iota
	T1
	T2
	T3
	T4
	T5
)

type Slot uint8

const (
	NoneSlot Slot = iota
	Weapon
	Chest
	Head
	Waist
	Foot
	Hand
	Neck
	Ring
)
