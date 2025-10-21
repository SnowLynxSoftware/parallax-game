package models

// WorldType represents the parallel world classification
type WorldType string

const (
	WorldTypeTutorial WorldType = "tutorial"
	WorldTypeFire     WorldType = "fire"
	WorldTypeIce      WorldType = "ice"
	WorldTypeTech     WorldType = "tech"
	WorldTypeNature   WorldType = "nature"
	WorldTypeVoid     WorldType = "void"
	WorldTypeLight    WorldType = "light"
)

// DifficultyLevel represents rift difficulty
type DifficultyLevel string

const (
	DifficultyTutorial  DifficultyLevel = "tutorial"
	DifficultyEasy      DifficultyLevel = "easy"
	DifficultyMedium    DifficultyLevel = "medium"
	DifficultyHard      DifficultyLevel = "hard"
	DifficultyLegendary DifficultyLevel = "legendary"
)

// ElementalType represents elemental affinities
type ElementalType string

const (
	ElementalTypeFire  ElementalType = "fire"
	ElementalTypeWater ElementalType = "water"
	ElementalTypeWind  ElementalType = "wind"
	ElementalTypeEarth ElementalType = "earth"
	ElementalTypeVoid  ElementalType = "void"
	ElementalTypeLight ElementalType = "light"
	ElementalTypeNone  ElementalType = "none"
)

// ItemRarity represents loot item rarity tiers
type ItemRarity string

const (
	ItemRarityCommon    ItemRarity = "common"
	ItemRarityUncommon  ItemRarity = "uncommon"
	ItemRarityRare      ItemRarity = "rare"
	ItemRarityEpic      ItemRarity = "epic"
	ItemRarityLegendary ItemRarity = "legendary"
)

// ItemType represents whether an item is equipment or consumable
type ItemType string

const (
	ItemTypeEquipment  ItemType = "equipment"
	ItemTypeConsumable ItemType = "consumable"
)

// EquipmentSlotType represents equipment slot categories
type EquipmentSlotType string

const (
	EquipmentSlotWeapon    EquipmentSlotType = "weapon"
	EquipmentSlotArmor     EquipmentSlotType = "armor"
	EquipmentSlotAccessory EquipmentSlotType = "accessory"
	EquipmentSlotArtifact  EquipmentSlotType = "artifact"
	EquipmentSlotRelic     EquipmentSlotType = "relic"
)
