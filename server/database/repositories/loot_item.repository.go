package repositories

import (
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type ILootItemRepository interface {
	GetAllLootItems() ([]*LootItemEntity, error)
	GetLootItemById(id int64) (*LootItemEntity, error)
	GetLootItemsByWorldType(worldType string) ([]*LootItemEntity, error)
	GetLootItemsByRarityAndWorldType(rarity string, worldType string) ([]*LootItemEntity, error)
}

type LootItemRepository struct {
	db *database.AppDataSource
}

func NewLootItemRepository(db *database.AppDataSource) ILootItemRepository {
	return &LootItemRepository{
		db: db,
	}
}

func (r *LootItemRepository) GetAllLootItems() ([]*LootItemEntity, error) {
	items := []*LootItemEntity{}
	sql := `SELECT * FROM loot_items WHERE is_archived = false ORDER BY rarity, name`
	err := r.db.DB.Select(&items, sql)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *LootItemRepository) GetLootItemById(id int64) (*LootItemEntity, error) {
	item := &LootItemEntity{}
	sql := `SELECT * FROM loot_items WHERE id = $1 AND is_archived = false`
	err := r.db.DB.Get(item, sql, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *LootItemRepository) GetLootItemsByWorldType(worldType string) ([]*LootItemEntity, error) {
	items := []*LootItemEntity{}
	sql := `SELECT * FROM loot_items WHERE world_type = $1 AND is_archived = false ORDER BY rarity, name`
	err := r.db.DB.Select(&items, sql, worldType)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *LootItemRepository) GetLootItemsByRarityAndWorldType(rarity string, worldType string) ([]*LootItemEntity, error) {
	items := []*LootItemEntity{}
	sql := `SELECT * FROM loot_items WHERE rarity = $1 AND world_type = $2 AND is_archived = false ORDER BY name`
	err := r.db.DB.Select(&items, sql, rarity, worldType)
	if err != nil {
		return nil, err
	}
	return items, nil
}
