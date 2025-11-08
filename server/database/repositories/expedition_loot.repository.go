package repositories

import (
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type IExpeditionLootRepository interface {
	CreateExpeditionLoot(expeditionId, lootItemId int64, quantity int) error
	GetLootByExpeditionId(expeditionId int64) ([]*ExpeditionLootEntity, error)
}

type ExpeditionLootRepository struct {
	db *database.AppDataSource
}

func NewExpeditionLootRepository(db *database.AppDataSource) IExpeditionLootRepository {
	return &ExpeditionLootRepository{
		db: db,
	}
}

func (r *ExpeditionLootRepository) CreateExpeditionLoot(expeditionId, lootItemId int64, quantity int) error {
	sql := `INSERT INTO expedition_loot (expedition_id, loot_item_id, quantity) VALUES ($1, $2, $3)`
	_, err := r.db.DB.Exec(sql, expeditionId, lootItemId, quantity)
	return err
}

func (r *ExpeditionLootRepository) GetLootByExpeditionId(expeditionId int64) ([]*ExpeditionLootEntity, error) {
	loot := []*ExpeditionLootEntity{}
	sql := `SELECT * FROM expedition_loot WHERE expedition_id = $1 AND is_archived = false ORDER BY created_at`
	err := r.db.DB.Select(&loot, sql, expeditionId)
	if err != nil {
		return nil, err
	}
	return loot, nil
}
