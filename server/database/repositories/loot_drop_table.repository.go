package repositories

import (
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type ILootDropTableRepository interface {
	GetDropTablesByRiftId(riftId int64) ([]*LootDropTableEntity, error)
}

type LootDropTableRepository struct {
	db *database.AppDataSource
}

func NewLootDropTableRepository(db *database.AppDataSource) ILootDropTableRepository {
	return &LootDropTableRepository{
		db: db,
	}
}

func (r *LootDropTableRepository) GetDropTablesByRiftId(riftId int64) ([]*LootDropTableEntity, error) {
	tables := []*LootDropTableEntity{}
	sql := `SELECT * FROM loot_drop_tables WHERE rift_id = $1 AND is_archived = false ORDER BY rarity`
	err := r.db.DB.Select(&tables, sql, riftId)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
