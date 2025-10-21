package repositories

import (
	"time"
)

// RiftEntity represents a parallel world rift that players can explore
type RiftEntity struct {
	ID                    int64      `json:"id" db:"id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt            *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived            bool       `json:"is_archived" db:"is_archived"`
	Name                  string     `json:"name" db:"name"`
	Description           string     `json:"description" db:"description"`
	WorldType             string     `json:"world_type" db:"world_type"`
	DurationMinutes       int        `json:"duration_minutes" db:"duration_minutes"`
	Difficulty            string     `json:"difficulty" db:"difficulty"`
	WeakToElement         string     `json:"weak_to_element" db:"weak_to_element"`
	UnlockRequirementText *string    `json:"unlock_requirement_text" db:"unlock_requirement_text"`
	Icon                  string     `json:"icon" db:"icon"`
}

// LootItemEntity represents a collectible item (equipment or consumable)
type LootItemEntity struct {
	ID                int64      `json:"id" db:"id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt        *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived        bool       `json:"is_archived" db:"is_archived"`
	Name              string     `json:"name" db:"name"`
	Description       string     `json:"description" db:"description"`
	Rarity            string     `json:"rarity" db:"rarity"`
	WorldType         string     `json:"world_type" db:"world_type"`
	ItemType          string     `json:"item_type" db:"item_type"`
	EquipmentSlot     *string    `json:"equipment_slot" db:"equipment_slot"`
	SpeedBonus        float64    `json:"speed_bonus" db:"speed_bonus"`
	LuckBonus         float64    `json:"luck_bonus" db:"luck_bonus"`
	PowerBonus        int        `json:"power_bonus" db:"power_bonus"`
	ElementalAffinity string     `json:"elemental_affinity" db:"elemental_affinity"`
	PowerValue        int        `json:"power_value" db:"power_value"`
	Icon              string     `json:"icon" db:"icon"`
}

// LootDropTableEntity represents drop rates for a rift/rarity combination
type LootDropTableEntity struct {
	ID              int64      `json:"id" db:"id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt      *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived      bool       `json:"is_archived" db:"is_archived"`
	RiftID          int64      `json:"rift_id" db:"rift_id"`
	Rarity          string     `json:"rarity" db:"rarity"`
	DropRatePercent float64    `json:"drop_rate_percent" db:"drop_rate_percent"`
	MinQuantity     int        `json:"min_quantity" db:"min_quantity"`
	MaxQuantity     int        `json:"max_quantity" db:"max_quantity"`
}

// TeamEntity represents a user's expedition team
type TeamEntity struct {
	ID                    int64      `json:"id" db:"id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt            *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived            bool       `json:"is_archived" db:"is_archived"`
	UserID                int64      `json:"user_id" db:"user_id"`
	TeamNumber            int        `json:"team_number" db:"team_number"`
	SpeedBonus            float64    `json:"speed_bonus" db:"speed_bonus"`
	LuckBonus             float64    `json:"luck_bonus" db:"luck_bonus"`
	PowerBonus            int        `json:"power_bonus" db:"power_bonus"`
	IsUnlocked            bool       `json:"is_unlocked" db:"is_unlocked"`
	EquippedWeaponSlot    *int64     `json:"equipped_weapon_slot" db:"equipped_weapon_slot"`
	EquippedArmorSlot     *int64     `json:"equipped_armor_slot" db:"equipped_armor_slot"`
	EquippedAccessorySlot *int64     `json:"equipped_accessory_slot" db:"equipped_accessory_slot"`
	EquippedArtifactSlot  *int64     `json:"equipped_artifact_slot" db:"equipped_artifact_slot"`
	EquippedRelicSlot     *int64     `json:"equipped_relic_slot" db:"equipped_relic_slot"`
}

// UserInventoryEntity represents a player's inventory item
type UserInventoryEntity struct {
	ID         int64      `json:"id" db:"id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived bool       `json:"is_archived" db:"is_archived"`
	UserID     int64      `json:"user_id" db:"user_id"`
	LootItemID int64      `json:"loot_item_id" db:"loot_item_id"`
	Quantity   int        `json:"quantity" db:"quantity"`
	AcquiredAt time.Time  `json:"acquired_at" db:"acquired_at"`
}

// ExpeditionEntity represents an expedition instance
type ExpeditionEntity struct {
	ID              int64      `json:"id" db:"id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt      *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived      bool       `json:"is_archived" db:"is_archived"`
	UserID          int64      `json:"user_id" db:"user_id"`
	TeamID          int64      `json:"team_id" db:"team_id"`
	RiftID          int64      `json:"rift_id" db:"rift_id"`
	StartTime       time.Time  `json:"start_time" db:"start_time"`
	DurationMinutes int        `json:"duration_minutes" db:"duration_minutes"`
	Completed       bool       `json:"completed" db:"completed"`
	Processed       bool       `json:"processed" db:"processed"`
	Claimed         bool       `json:"claimed" db:"claimed"`
}

// ExpeditionLootEntity represents loot audit trail for an expedition
type ExpeditionLootEntity struct {
	ID           int64      `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt   *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived   bool       `json:"is_archived" db:"is_archived"`
	ExpeditionID int64      `json:"expedition_id" db:"expedition_id"`
	LootItemID   int64      `json:"loot_item_id" db:"loot_item_id"`
	Quantity     int        `json:"quantity" db:"quantity"`
}
