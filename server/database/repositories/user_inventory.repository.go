package repositories

import (
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type IUserInventoryRepository interface {
	GetInventoryByUserId(userId int64) ([]*UserInventoryEntity, error)
	GetInventoryById(inventoryId int64) (*UserInventoryEntity, error)
	GetEquipmentByUserId(userId int64) ([]*UserInventoryEntity, error)
	GetConsumablesByUserId(userId int64) ([]*UserInventoryEntity, error)
	AddLoot(userId int64, lootItemId int64, itemType string) (*UserInventoryEntity, error)
	ConsumeLoot(inventoryId int64) error
	GetInventoryByUserAndItem(userId int64, lootItemId int64) (*UserInventoryEntity, error)
	HasItemByName(userId int64, itemName string) (bool, error)
}

type UserInventoryRepository struct {
	db *database.AppDataSource
}

func NewUserInventoryRepository(db *database.AppDataSource) IUserInventoryRepository {
	return &UserInventoryRepository{
		db: db,
	}
}

func (r *UserInventoryRepository) GetInventoryByUserId(userId int64) ([]*UserInventoryEntity, error) {
	inventory := []*UserInventoryEntity{}
	sql := `SELECT * FROM user_inventory WHERE user_id = $1 AND is_archived = false ORDER BY acquired_at DESC`
	err := r.db.DB.Select(&inventory, sql, userId)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *UserInventoryRepository) GetInventoryById(inventoryId int64) (*UserInventoryEntity, error) {
	item := &UserInventoryEntity{}
	sql := `SELECT * FROM user_inventory WHERE id = $1 AND is_archived = false`
	err := r.db.DB.Get(item, sql, inventoryId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *UserInventoryRepository) GetEquipmentByUserId(userId int64) ([]*UserInventoryEntity, error) {
	inventory := []*UserInventoryEntity{}
	sql := `SELECT ui.* FROM user_inventory ui
			JOIN loot_items li ON li.id = ui.loot_item_id
			WHERE ui.user_id = $1 AND li.item_type = 'equipment' AND ui.is_archived = false
			ORDER BY ui.acquired_at DESC`
	err := r.db.DB.Select(&inventory, sql, userId)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *UserInventoryRepository) GetConsumablesByUserId(userId int64) ([]*UserInventoryEntity, error) {
	inventory := []*UserInventoryEntity{}
	sql := `SELECT ui.* FROM user_inventory ui
			JOIN loot_items li ON li.id = ui.loot_item_id
			WHERE ui.user_id = $1 AND li.item_type = 'consumable' AND ui.is_archived = false
			ORDER BY li.name`
	err := r.db.DB.Select(&inventory, sql, userId)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *UserInventoryRepository) GetInventoryByUserAndItem(userId int64, lootItemId int64) (*UserInventoryEntity, error) {
	item := &UserInventoryEntity{}
	sql := `SELECT * FROM user_inventory WHERE user_id = $1 AND loot_item_id = $2 AND is_archived = false LIMIT 1`
	err := r.db.DB.Get(item, sql, userId, lootItemId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

// AddLoot handles both equipment (always new row) and consumable (upsert) logic
func (r *UserInventoryRepository) AddLoot(userId int64, lootItemId int64, itemType string) (*UserInventoryEntity, error) {
	if itemType == "equipment" {
		// Equipment: Always insert new row with quantity = 1
		sql := `INSERT INTO user_inventory (user_id, loot_item_id, quantity, acquired_at)
				VALUES ($1, $2, 1, NOW())
				RETURNING id, created_at, modified_at, is_archived, user_id, loot_item_id, quantity, acquired_at`
		item := &UserInventoryEntity{}
		err := r.db.DB.QueryRowx(sql, userId, lootItemId).StructScan(item)
		if err != nil {
			return nil, err
		}
		return item, nil
	} else {
		// Consumable: Upsert (insert or increment quantity)
		sql := `INSERT INTO user_inventory (user_id, loot_item_id, quantity, acquired_at)
				VALUES ($1, $2, 1, NOW())
				ON CONFLICT ON CONSTRAINT user_inventory_pkey DO NOTHING
				RETURNING id, created_at, modified_at, is_archived, user_id, loot_item_id, quantity, acquired_at`

		item := &UserInventoryEntity{}
		err := r.db.DB.QueryRowx(sql, userId, lootItemId).StructScan(item)

		if err != nil {
			// If conflict, increment existing
			existing, err := r.GetInventoryByUserAndItem(userId, lootItemId)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				sql := `UPDATE user_inventory SET quantity = quantity + 1, modified_at = NOW() WHERE id = $1
						RETURNING id, created_at, modified_at, is_archived, user_id, loot_item_id, quantity, acquired_at`
				err = r.db.DB.QueryRowx(sql, existing.ID).StructScan(item)
				if err != nil {
					return nil, err
				}
				return item, nil
			}
			return nil, err
		}

		return item, nil
	}
}

// ConsumeLoot decrements quantity and archives item if quantity reaches 0
// Uses a single atomic UPDATE with CASE to avoid CHECK constraint violations
// Since quantity has CHECK (quantity >= 1), we must NOT let it go below 1
// When quantity=1, we archive it. When quantity>1, we decrement.
func (r *UserInventoryRepository) ConsumeLoot(inventoryId int64) error {
	sql := `UPDATE user_inventory 
			SET quantity = CASE 
				WHEN quantity > 1 THEN quantity - 1
				ELSE 1
			END,
			is_archived = CASE
				WHEN quantity = 1 THEN true
				ELSE false
			END,
			modified_at = NOW()
			WHERE id = $1`

	_, err := r.db.DB.Exec(sql, inventoryId)
	return err
}

// HasItemByName checks if a user has a specific item by its name
func (r *UserInventoryRepository) HasItemByName(userId int64, itemName string) (bool, error) {
	var count int
	sql := `SELECT COUNT(*) FROM user_inventory ui
			JOIN loot_items li ON li.id = ui.loot_item_id
			WHERE ui.user_id = $1 AND li.name = $2 AND ui.is_archived = false AND li.is_archived = false`
	err := r.db.DB.Get(&count, sql, userId, itemName)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
