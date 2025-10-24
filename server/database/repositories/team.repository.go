package repositories

import (
	"database/sql"
	"fmt"

	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type ITeamRepository interface {
	GetTeamsByUserId(userId int64) ([]*TeamEntity, error)
	GetTeamById(teamId int64) (*TeamEntity, error)
	CreateTeamsForUser(userId int64) error
	UpdateTeamStats(teamId int64, speedBonus, luckBonus float64, powerBonus int) error
	EquipItem(teamId int64, slot string, inventoryId *int64) error
	UnequipItem(teamId int64, slot string) error
	UnlockTeam(teamId int64) error
	GetTeamsByUserIdWithSlot(userId int64, inventoryId int64) (*TeamEntity, *string, error)
}

type TeamRepository struct {
	db *database.AppDataSource
}

func NewTeamRepository(db *database.AppDataSource) ITeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) GetTeamsByUserId(userId int64) ([]*TeamEntity, error) {
	teams := []*TeamEntity{}
	sql := `SELECT * FROM teams WHERE user_id = $1 AND is_archived = false ORDER BY team_number`
	err := r.db.DB.Select(&teams, sql, userId)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) GetTeamById(teamId int64) (*TeamEntity, error) {
	team := &TeamEntity{}
	sql := `SELECT * FROM teams WHERE id = $1 AND is_archived = false`
	err := r.db.DB.Get(team, sql, teamId)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepository) CreateTeamsForUser(userId int64) error {
	// Create Team 1 (unlocked)
	sql := `INSERT INTO teams (user_id, team_number, is_unlocked) VALUES ($1, 1, true)`
	_, err := r.db.DB.Exec(sql, userId)
	if err != nil {
		return err
	}

	// Create Teams 2-5 (locked)
	for i := 2; i <= 5; i++ {
		sql := `INSERT INTO teams (user_id, team_number, is_unlocked) VALUES ($1, $2, false)`
		_, err := r.db.DB.Exec(sql, userId, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) UpdateTeamStats(teamId int64, speedBonus, luckBonus float64, powerBonus int) error {
	sql := `UPDATE teams 
			SET speed_bonus = speed_bonus + $1, 
				luck_bonus = luck_bonus + $2, 
				power_bonus = power_bonus + $3,
				modified_at = NOW()
			WHERE id = $4`
	_, err := r.db.DB.Exec(sql, speedBonus, luckBonus, powerBonus, teamId)
	return err
}

func (r *TeamRepository) EquipItem(teamId int64, slot string, inventoryId *int64) error {
	var sql string
	switch slot {
	case "weapon":
		sql = `UPDATE teams SET equipped_weapon_slot = $1, modified_at = NOW() WHERE id = $2`
	case "armor":
		sql = `UPDATE teams SET equipped_armor_slot = $1, modified_at = NOW() WHERE id = $2`
	case "accessory":
		sql = `UPDATE teams SET equipped_accessory_slot = $1, modified_at = NOW() WHERE id = $2`
	case "artifact":
		sql = `UPDATE teams SET equipped_artifact_slot = $1, modified_at = NOW() WHERE id = $2`
	case "relic":
		sql = `UPDATE teams SET equipped_relic_slot = $1, modified_at = NOW() WHERE id = $2`
	default:
		return fmt.Errorf("invalid equipment slot: %s", slot)
	}

	_, err := r.db.DB.Exec(sql, inventoryId, teamId)
	return err
}

func (r *TeamRepository) UnequipItem(teamId int64, slot string) error {
	return r.EquipItem(teamId, slot, nil)
}

func (r *TeamRepository) UnlockTeam(teamId int64) error {
	sql := `UPDATE teams SET is_unlocked = true, modified_at = NOW() WHERE id = $1`
	_, err := r.db.DB.Exec(sql, teamId)
	return err
}

// GetTeamsByUserIdWithSlot checks if an inventory item is equipped to any of user's teams
// Returns the team, slot name, and error (if any)
func (r *TeamRepository) GetTeamsByUserIdWithSlot(userId int64, inventoryId int64) (*TeamEntity, *string, error) {
	team := &TeamEntity{}
	var slot sql.NullString

	sqlQuery := `SELECT t.id, t.created_at, t.modified_at, t.is_archived, t.user_id, t.team_number,
			t.speed_bonus, t.luck_bonus, t.power_bonus, t.is_unlocked,
			t.equipped_weapon_slot, t.equipped_armor_slot, t.equipped_accessory_slot,
			t.equipped_artifact_slot, t.equipped_relic_slot,
			CASE
				WHEN t.equipped_weapon_slot = $2 THEN 'weapon'
				WHEN t.equipped_armor_slot = $2 THEN 'armor'
				WHEN t.equipped_accessory_slot = $2 THEN 'accessory'
				WHEN t.equipped_artifact_slot = $2 THEN 'artifact'
				WHEN t.equipped_relic_slot = $2 THEN 'relic'
			END as slot
			FROM teams t
			WHERE t.user_id = $1 AND (
				t.equipped_weapon_slot = $2 OR
				t.equipped_armor_slot = $2 OR
				t.equipped_accessory_slot = $2 OR
				t.equipped_artifact_slot = $2 OR
				t.equipped_relic_slot = $2
			) AND t.is_archived = false
			LIMIT 1`

	err := r.db.DB.QueryRowx(sqlQuery, userId, inventoryId).Scan(
		&team.ID, &team.CreatedAt, &team.ModifiedAt, &team.IsArchived, &team.UserID, &team.TeamNumber,
		&team.SpeedBonus, &team.LuckBonus, &team.PowerBonus, &team.IsUnlocked,
		&team.EquippedWeaponSlot, &team.EquippedArmorSlot, &team.EquippedAccessorySlot,
		&team.EquippedArtifactSlot, &team.EquippedRelicSlot,
		&slot,
	)

	if err == sql.ErrNoRows {
		return nil, nil, nil // Not equipped anywhere
	}

	if err != nil {
		return nil, nil, err
	}

	if slot.Valid {
		slotStr := slot.String
		return team, &slotStr, nil
	}

	return nil, nil, nil // Not equipped anywhere
}
