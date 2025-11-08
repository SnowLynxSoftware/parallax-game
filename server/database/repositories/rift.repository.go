package repositories

import (
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type IRiftRepository interface {
	GetAllRifts() ([]*RiftEntity, error)
	GetRiftById(id int64) (*RiftEntity, error)
	GetRiftsByDifficulty(difficulty string) ([]*RiftEntity, error)
}

type RiftRepository struct {
	db *database.AppDataSource
}

func NewRiftRepository(db *database.AppDataSource) IRiftRepository {
	return &RiftRepository{
		db: db,
	}
}

func (r *RiftRepository) GetAllRifts() ([]*RiftEntity, error) {
	rifts := []*RiftEntity{}
	sql := `SELECT * FROM rifts WHERE is_archived = false ORDER BY difficulty, duration_minutes`
	err := r.db.DB.Select(&rifts, sql)
	if err != nil {
		return nil, err
	}
	return rifts, nil
}

func (r *RiftRepository) GetRiftById(id int64) (*RiftEntity, error) {
	rift := &RiftEntity{}
	sql := `SELECT * FROM rifts WHERE id = $1 AND is_archived = false`
	err := r.db.DB.Get(rift, sql, id)
	if err != nil {
		return nil, err
	}
	return rift, nil
}

func (r *RiftRepository) GetRiftsByDifficulty(difficulty string) ([]*RiftEntity, error) {
	rifts := []*RiftEntity{}
	sql := `SELECT * FROM rifts WHERE difficulty = $1 AND is_archived = false ORDER BY duration_minutes`
	err := r.db.DB.Select(&rifts, sql, difficulty)
	if err != nil {
		return nil, err
	}
	return rifts, nil
}
