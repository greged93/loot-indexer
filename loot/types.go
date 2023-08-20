package loot

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/NethermindEth/juno/core/felt"
)

func FeltToBigInt(f *felt.Felt) (*big.Int, bool) {
	return new(big.Int).SetString(f.String(), 0)
}

type SqlBigInt big.Int

func (sqlB *SqlBigInt) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	result := new(big.Int).SetBytes(bytes)
	*sqlB = SqlBigInt(*result)
	return nil
}

func (sqlB SqlBigInt) Value() (driver.Value, error) {
	b := big.Int(sqlB)
	return json.Marshal(b.Bytes())
}

type Packable interface{}

func Pack[T Packable](s T) uint64 {
	val := reflect.ValueOf(s)
	var acc uint64 = 0
	for i := 0; i < val.NumField(); i++ {
		acc |= val.Field(i).Uint() << uint(8*i)
	}
	return acc
}

func Unpack[T Packable](packed uint64, s *T) error {
	if s == nil {
		return fmt.Errorf("cannot unpack into nil Packable")
	}

	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		val.Field(i).SetUint(uint64(0xff & (packed >> uint(8*i))))
	}
	v, ok := val.Interface().(T)
	if !ok {
		return fmt.Errorf("expected %T, got %T", s, val.Interface())
	}
	*s = v
	return nil
}

type Adventurer struct {
	AdventurerID        SqlBigInt
	Level               uint16
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
	Name      SqlBigInt
	HomeRealm uint16
	Class     uint8
	Entropy   SqlBigInt
}

type ItemPrimitive struct {
	Id       uint8  // 7 bits
	Xp       uint16 // 9 bits
	Metadata uint8  // 5 bits
}

func (i *ItemPrimitive) Scan(value interface{}) error {
	v, ok := value.(uint32)
	if !ok {
		return fmt.Errorf("expected uint32, got %T", value)
	}
	return Unpack(uint64(v), i)
}

func (i ItemPrimitive) Value() (driver.Value, error) {
	return uint32(Pack(i)), nil
}

type Stats struct {
	Strength     uint8 // 5 bits
	Dexterity    uint8 // 5 bits
	Vitality     uint8 // 5 bits
	Intelligence uint8 // 5 bits
	Wisdom       uint8 // 5 bits
	Charisma     uint8 // 5 bits
}

func (s *Stats) Scan(value interface{}) error {
	v, ok := value.(uint64)
	if !ok {
		return fmt.Errorf("expected uint64, got %T", value)
	}
	return Unpack(v, s)
}

func (s Stats) Value() (driver.Value, error) {
	return Pack(s), nil
}

type AdventurerState struct {
	Owner        SqlBigInt
	AdventurerId SqlBigInt
	Adventurer
}

type AdventurerStateWithBag struct {
	AdventurerState
	Bag
}

func (a *AdventurerStateWithBag) ApplyEvent(event []byte) error {
	var rawEvent RawEvent
	err := json.Unmarshal(event, &rawEvent)
	if err != nil {
		return err
	}

	switch rawEvent.Type {
	case StarGameType:
		var startGame StartGame
		err := json.Unmarshal(rawEvent.Event, &startGame)
		if err != nil {
			return err
		}
		a.AdventurerState = startGame.AdventurerState
	case DiscoveredHealthType:
		var discoveredHealth DiscoveredHealth
		err := json.Unmarshal(rawEvent.Event, &discoveredHealth)
		if err != nil {
			return err
		}
		a.Health += discoveredHealth.HealthAmount
	case DiscoveredGoldType:
		var discoveredGold DiscoveredGold
		err := json.Unmarshal(rawEvent.Event, &discoveredGold)
		if err != nil {
			return err
		}
		a.Gold += discoveredGold.Goldamount
	case DiscoveredXPType:
		var discoveredXP DiscoveredXP
		err := json.Unmarshal(rawEvent.Event, &discoveredXP)
		if err != nil {
			return err
		}
		a.Xp += discoveredXP.XpAmount
	case DodgedObstacleType:
		var dodgedObstacle DodgedObstacle
		err := json.Unmarshal(rawEvent.Event, &dodgedObstacle)
		if err != nil {
			return err
		}
		a.Health -= dodgedObstacle.DamageTaken
		a.Xp += dodgedObstacle.XpEarnedAdventurer
		a.AddXpWeapons(dodgedObstacle.XpEarnedItems)
	case HitByObstacleType:
		var hitByObstacle HitByObstacle
		err := json.Unmarshal(rawEvent.Event, &hitByObstacle)
		if err != nil {
			return err
		}
		a.Health -= hitByObstacle.DamageTaken
		a.Xp += hitByObstacle.XpEarnedAdventurer
		a.AddXpWeapons(hitByObstacle.XpEarnedItems)

	case DiscoveredBeastType:
		var discoveredBeast DiscoveredBeast
		err := json.Unmarshal(rawEvent.Event, &discoveredBeast)
		if err != nil {
			return err
		}
		a.BeastHealth = discoveredBeast.BeastHealth
	case AmbushedByBeastType:
		var ambushedByBeast AmbushedByBeast
		err := json.Unmarshal(rawEvent.Event, &ambushedByBeast)
		if err != nil {
			return err
		}
		a.BeastHealth = ambushedByBeast.BeastHealth
		a.Health -= ambushedByBeast.Damage
	case AttackedBeastType:
		var attackedBeast AttackedBeast
		err := json.Unmarshal(rawEvent.Event, &attackedBeast)
		if err != nil {
			return err
		}
		a.BeastHealth -= attackedBeast.Damage
	case AttackedByBeastType:
		var attackedByBeast AttackedByBeast
		err := json.Unmarshal(rawEvent.Event, &attackedByBeast)
		if err != nil {
			return err
		}
		a.Health -= attackedByBeast.Damage
	case SlayedBeastType:
		var slayedBeast SlayedBeast
		err := json.Unmarshal(rawEvent.Event, &slayedBeast)
		if err != nil {
			return err
		}
		a.Xp += slayedBeast.XpEarnedAdventurer
		a.Gold += slayedBeast.GoldEarned
		a.BeastHealth = 0
		a.AddXpWeapons(slayedBeast.XpEarnedItems)

	case FleeFailedType:
		break
	case FleeSucceededType:
		var fleeSucceeded FleeSucceeded
		err := json.Unmarshal(rawEvent.Event, &fleeSucceeded)
		if err != nil {
			return err
		}
		a.BeastHealth = 0
	case PurchasedPotionsType:
		var purchasedPotions PurchasedPotions
		err := json.Unmarshal(rawEvent.Event, &purchasedPotions)
		if err != nil {
			return err
		}
		a.Gold -= purchasedPotions.Cost
		a.Health += uint16(purchasedPotions.Quantity) * purchasedPotions.Health
	case PurchasedItemsType:
		var purchasedItems PurchasedItems
		err := json.Unmarshal(rawEvent.Event, &purchasedItems)
		if err != nil {
			return err
		}
		for _, loot := range purchasedItems.Purchases {
			a.Gold -= loot.Price
			a.Bag.addItem(loot.Loot)
		}
	case EquippedItemsType:
		var equippedItems EquippedItems
		err := json.Unmarshal(rawEvent.Event, &equippedItems)
		if err != nil {
			return err
		}
		// TODO: equip items
	case DroppedItemsType:
		var droppedItems DroppedItems
		err := json.Unmarshal(rawEvent.Event, &droppedItems)
		if err != nil {
			return err
		}
		// TODO drop items
	case ItemLeveledUpType:
		var itemLeveledUp ItemLeveledUp
		err := json.Unmarshal(rawEvent.Event, &itemLeveledUp)
		if err != nil {
			return err
		}
		// TODO level up items
	case ItemSpecialUnlockedType:
		var itemSpecialUnlocked ItemSpecialUnlocked
		err := json.Unmarshal(rawEvent.Event, &itemSpecialUnlocked)
		if err != nil {
			return err
		}
		// TODO unlock item special
	case NewHighScoreType:
		var newHighScore NewHighScore
		err := json.Unmarshal(rawEvent.Event, &newHighScore)
		if err != nil {
			return err
		}
		// TODO update high score
	case AdventurerDiedType:
		var adventurerDied AdventurerDied
		err := json.Unmarshal(rawEvent.Event, &adventurerDied)
		if err != nil {
			return err
		}
		a.Health = 0
	case ShopAvailableType:
		var shopAvailable ShopAvailable
		err := json.Unmarshal(rawEvent.Event, &shopAvailable)
		if err != nil {
			return err
		}
		// TODO update shop
	case AdventurerLeveledUpType:
		var adventurerLeveledUp AdventurerLeveledUp
		err := json.Unmarshal(rawEvent.Event, &adventurerLeveledUp)
		if err != nil {
			return err
		}
		a.Level += 1
	case NewItemsAvailableType:
		var newItemsAvailable NewItemsAvailable
		err := json.Unmarshal(rawEvent.Event, &newItemsAvailable)
		if err != nil {
			return err
		}
		// TODO update shop
	case IdleDamagePenaltyType:
		var idleDamagePenalty IdleDamagePenalty
		err := json.Unmarshal(rawEvent.Event, &idleDamagePenalty)
		if err != nil {
			return err
		}
		a.Health -= idleDamagePenalty.DamageTaken
	case UpgradeAvailableType:
		var upgradeAvailable UpgradeAvailable
		err := json.Unmarshal(rawEvent.Event, &upgradeAvailable)
		if err != nil {
			return err
		}
		// TODO upgrade
	case AdventurerUpgradedType:
		var adventurerUpgraded AdventurerUpgraded
		err := json.Unmarshal(rawEvent.Event, &adventurerUpgraded)
		if err != nil {
			return err
		}
		a.Upgrade(adventurerUpgraded)
	default:
		return fmt.Errorf("unknown event type %d", rawEvent.Type)
	}
	return nil
}

func (a *AdventurerState) AddXpWeapons(xp uint16) {
	a.Weapon.Xp += xp
	a.Chest.Xp += xp
	a.Head.Xp += xp
	a.Waist.Xp += xp
	a.Foot.Xp += xp
	a.Hand.Xp += xp
	a.Neck.Xp += xp
	a.Ring.Xp += xp
}

func (a *AdventurerState) Upgrade(upgrade AdventurerUpgraded) {
	a.Stats.Strength += upgrade.StrengthIncrease
	a.Stats.Dexterity += upgrade.DexterityIncrease
	a.Stats.Vitality += upgrade.VitalityIncrease
	a.Stats.Intelligence += upgrade.IntelligenceIncrease
	a.Stats.Wisdom += upgrade.WisdomIncrease
	a.Stats.Charisma += upgrade.CharismaIncrease
}

type ItemPrimitives [11]ItemPrimitive

func (items *ItemPrimitives) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(val, items)
}

func (items ItemPrimitives) Value() (driver.Value, error) {
	return json.Marshal(items)
}

type Bag struct {
	Items   ItemPrimitives
	Mutated bool
}

func (b *Bag) addItem(loot Loot) {
	newItem := ItemPrimitive{
		Id:       loot.Id,
		Xp:       0,
		Metadata: 0, // TODO
	}
	for i, item := range b.Items {
		if item.Id == 0 {
			b.Items[i] = newItem
			return
		}
	}
}

type CombatSpec struct {
	Tier     Tier
	ItemType Type
	Level    uint16
	Specials uint32 // packed bits of SpecialPowers
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
