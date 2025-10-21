package models

// Request DTOs

type StartExpeditionDTO struct {
	TeamID int64 `json:"team_id"`
	RiftID int64 `json:"rift_id"`
}

type EquipItemDTO struct {
	TeamID      int64  `json:"team_id"`
	Slot        string `json:"slot"`
	InventoryID int64  `json:"inventory_id"`
}

type UnequipItemDTO struct {
	TeamID int64  `json:"team_id"`
	Slot   string `json:"slot"`
}

type ConsumeItemDTO struct {
	TeamID      int64 `json:"team_id"`
	InventoryID int64 `json:"inventory_id"`
}

type UnlockTeamDTO struct {
	TeamID int64 `json:"team_id"`
}

// Response DTOs

type TeamStatsDTO struct {
	Speed float64 `json:"speed"`
	Luck  float64 `json:"luck"`
	Power int     `json:"power"`
}

type EquippedItemDTO struct {
	InventoryID       *int64  `json:"inventory_id"`
	LootItemID        *int64  `json:"loot_item_id"`
	Name              *string `json:"name"`
	Icon              *string `json:"icon"`
	Rarity            *string `json:"rarity"`
	SpeedBonus        float64 `json:"speed_bonus"`
	LuckBonus         float64 `json:"luck_bonus"`
	PowerBonus        int     `json:"power_bonus"`
	ElementalAffinity *string `json:"elemental_affinity"`
}

type TeamResponseDTO struct {
	ID                int64            `json:"id"`
	TeamNumber        int              `json:"team_number"`
	IsUnlocked        bool             `json:"is_unlocked"`
	BaseStats         TeamStatsDTO     `json:"base_stats"`
	TotalStats        TeamStatsDTO     `json:"total_stats"`
	EquippedWeapon    *EquippedItemDTO `json:"equipped_weapon"`
	EquippedArmor     *EquippedItemDTO `json:"equipped_armor"`
	EquippedAccessory *EquippedItemDTO `json:"equipped_accessory"`
	EquippedArtifact  *EquippedItemDTO `json:"equipped_artifact"`
	EquippedRelic     *EquippedItemDTO `json:"equipped_relic"`
}

type LootItemResponseDTO struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Rarity            string  `json:"rarity"`
	WorldType         string  `json:"world_type"`
	ItemType          string  `json:"item_type"`
	EquipmentSlot     *string `json:"equipment_slot"`
	SpeedBonus        float64 `json:"speed_bonus"`
	LuckBonus         float64 `json:"luck_bonus"`
	PowerBonus        int     `json:"power_bonus"`
	ElementalAffinity string  `json:"elemental_affinity"`
	PowerValue        int     `json:"power_value"`
	Icon              string  `json:"icon"`
}

type InventoryItemResponseDTO struct {
	InventoryID int64               `json:"inventory_id"`
	Quantity    int                 `json:"quantity"`
	AcquiredAt  string              `json:"acquired_at"`
	IsEquipped  bool                `json:"is_equipped"`
	LootItem    LootItemResponseDTO `json:"loot_item"`
}

type InventoryResponseDTO struct {
	Equipment   []InventoryItemResponseDTO `json:"equipment"`
	Consumables []InventoryItemResponseDTO `json:"consumables"`
}

type RiftResponseDTO struct {
	ID                    int64   `json:"id"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	WorldType             string  `json:"world_type"`
	DurationMinutes       int     `json:"duration_minutes"`
	Difficulty            string  `json:"difficulty"`
	WeakToElement         string  `json:"weak_to_element"`
	UnlockRequirementText *string `json:"unlock_requirement_text"`
	Icon                  string  `json:"icon"`
	IsUnlocked            bool    `json:"is_unlocked"`
}

type ExpeditionResponseDTO struct {
	ID              int64                  `json:"id"`
	TeamID          int64                  `json:"team_id"`
	TeamNumber      int                    `json:"team_number"`
	RiftID          int64                  `json:"rift_id"`
	RiftName        string                 `json:"rift_name"`
	StartTime       string                 `json:"start_time"`
	DurationMinutes int                    `json:"duration_minutes"`
	CompletionTime  string                 `json:"completion_time"`
	TimeRemaining   *int                   `json:"time_remaining"` // seconds, null if completed
	IsCompleted     bool                   `json:"is_completed"`
	IsClaimed       bool                   `json:"is_claimed"`
	Loot            *[]LootItemResponseDTO `json:"loot,omitempty"` // Only if claimed
}

type ExpeditionRewardsDTO struct {
	ExpeditionID int64                 `json:"expedition_id"`
	Loot         []LootItemResponseDTO `json:"loot"`
}
