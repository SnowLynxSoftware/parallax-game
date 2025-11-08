package repositories

import (
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type IExpeditionRepository interface {
	CreateExpedition(userId, teamId, riftId int64, durationMinutes int) (*ExpeditionEntity, error)
	GetExpeditionById(expeditionId int64) (*ExpeditionEntity, error)
	GetActiveExpeditionsByUserId(userId int64) ([]*ExpeditionEntity, error)
	GetCompletedExpeditionsByUserId(userId int64, limit int) ([]*ExpeditionEntity, error)
	GetCompletedExpeditionsCount(userId int64) (int, error)
	MarkCompleted(expeditionId int64) error
	MarkProcessed(expeditionId int64) error
	MarkClaimed(expeditionId int64) error
}

type ExpeditionRepository struct {
	db *database.AppDataSource
}

func NewExpeditionRepository(db *database.AppDataSource) IExpeditionRepository {
	return &ExpeditionRepository{
		db: db,
	}
}

func (r *ExpeditionRepository) CreateExpedition(userId, teamId, riftId int64, durationMinutes int) (*ExpeditionEntity, error) {
	expedition := &ExpeditionEntity{}
	sql := `INSERT INTO expeditions (user_id, team_id, rift_id, start_time, duration_minutes, completed, processed, claimed)
			VALUES ($1, $2, $3, NOW(), $4, false, false, false)
			RETURNING id, created_at, modified_at, is_archived, user_id, team_id, rift_id, start_time, duration_minutes, completed, processed, claimed`
	err := r.db.DB.QueryRowx(sql, userId, teamId, riftId, durationMinutes).StructScan(expedition)
	if err != nil {
		return nil, err
	}
	return expedition, nil
}

func (r *ExpeditionRepository) GetExpeditionById(expeditionId int64) (*ExpeditionEntity, error) {
	expedition := &ExpeditionEntity{}
	sql := `SELECT * FROM expeditions WHERE id = $1 AND is_archived = false`
	err := r.db.DB.Get(expedition, sql, expeditionId)
	if err != nil {
		return nil, err
	}
	return expedition, nil
}

func (r *ExpeditionRepository) GetActiveExpeditionsByUserId(userId int64) ([]*ExpeditionEntity, error) {
	expeditions := []*ExpeditionEntity{}
	sql := `SELECT * FROM expeditions 
			WHERE user_id = $1 AND completed = false AND is_archived = false
			ORDER BY start_time DESC`
	err := r.db.DB.Select(&expeditions, sql, userId)
	if err != nil {
		return nil, err
	}
	return expeditions, nil
}

func (r *ExpeditionRepository) GetCompletedExpeditionsByUserId(userId int64, limit int) ([]*ExpeditionEntity, error) {
	expeditions := []*ExpeditionEntity{}
	sql := `SELECT * FROM expeditions 
			WHERE user_id = $1 AND completed = true AND is_archived = false
			ORDER BY start_time DESC
			LIMIT $2`
	err := r.db.DB.Select(&expeditions, sql, userId, limit)
	if err != nil {
		return nil, err
	}
	return expeditions, nil
}

func (r *ExpeditionRepository) GetCompletedExpeditionsCount(userId int64) (int, error) {
	var count int
	sql := `SELECT COUNT(*) FROM expeditions WHERE user_id = $1 AND completed = true AND is_archived = false`
	err := r.db.DB.Get(&count, sql, userId)
	return count, err
}

func (r *ExpeditionRepository) MarkCompleted(expeditionId int64) error {
	sql := `UPDATE expeditions SET completed = true, modified_at = NOW() WHERE id = $1`
	_, err := r.db.DB.Exec(sql, expeditionId)
	return err
}

func (r *ExpeditionRepository) MarkProcessed(expeditionId int64) error {
	sql := `UPDATE expeditions SET processed = true, modified_at = NOW() WHERE id = $1`
	_, err := r.db.DB.Exec(sql, expeditionId)
	return err
}

func (r *ExpeditionRepository) MarkClaimed(expeditionId int64) error {
	// Check if expedition is completed
	expedition, err := r.GetExpeditionById(expeditionId)
	if err != nil {
		return err
	}

	// Mark as completed if not already (in case time has passed)
	completionTime := expedition.StartTime.Add(time.Duration(expedition.DurationMinutes) * time.Minute)
	if time.Now().After(completionTime) && !expedition.Completed {
		err = r.MarkCompleted(expeditionId)
		if err != nil {
			return err
		}
	}

	sql := `UPDATE expeditions SET claimed = true, modified_at = NOW() WHERE id = $1`
	_, err = r.db.DB.Exec(sql, expeditionId)
	return err
}
