package repositories

import (
	"fmt"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type UserEntity struct {
	ID           int64      `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt   *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived   bool       `json:"is_archived" db:"is_archived"`
	Email        string     `json:"email" db:"email"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	AvatarURL    *string    `json:"avatar_url" db:"avatar_url"`
	ProfileText  *string    `json:"profile_text" db:"profile_text"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	PasswordHash *string    `json:"-" db:"password_hash"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
}

type IUserRepository interface {
	GetUsersCount(searchString string, statusFilter string, userTypeFilter string) (*int, error)
	GetUsers(pageSize int, offset int, searchString string, statusFilter string, userTypeFilter string) ([]*UserEntity, error)
	GetUserById(id int) (*UserEntity, error)
	GetUserByEmail(email string) (*UserEntity, error)
	CreateNewUser(dto *models.UserCreateDTO) (*UserEntity, error)
	MarkUserVerified(userId *int) (bool, error)
	UpdateUser(dto *models.UserUpdateDTO, userId *int) (*UserEntity, error)
	UpdateUserLastLogin(userId *int) (bool, error)
	UpdateUserPassword(userId *int, password string) (bool, error)
	ToggleUserArchived(userId *int) error
}

type UserRepository struct {
	db *database.AppDataSource
}

func NewUserRepository(db *database.AppDataSource) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

// ToggleUserArchived toggles the archived status of a user and clears the password hash.
func (r *UserRepository) ToggleUserArchived(userId *int) error {
	sql := `UPDATE users SET is_archived = NOT is_archived, password_hash = '', modified_at = NOW() WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsersCount(searchString string, statusFilter string, userTypeFilter string) (*int, error) {
	count := new(int)
	sql := `SELECT
		COUNT(*) as count
	FROM users
	WHERE ((email LIKE '%' || $1 || '%') OR (display_name LIKE '%' || $1 || '%'))`

	// Build dynamic WHERE clause for filtering
	args := []interface{}{searchString}
	argIndex := 2

	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false AND is_banned = false`
		case "archived":
			sql += ` AND is_archived = true`
		case "banned":
			sql += ` AND is_banned = true`
		}
	}

	if userTypeFilter != "" {
		sql += ` AND user_type_key = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, userTypeFilter)
		argIndex++
	}

	err := r.db.DB.Get(&count, sql, args...)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (r *UserRepository) GetUsers(pageSize int, offset int, searchString string, statusFilter string, userTypeFilter string) ([]*UserEntity, error) {
	users := []*UserEntity{}
	sql := `SELECT
		*
	FROM users
	WHERE ((email LIKE '%' || $3 || '%') OR (display_name LIKE '%' || $3 || '%'))`

	// Build dynamic WHERE clause for filtering
	args := []interface{}{pageSize, offset, searchString}
	argIndex := 4

	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false AND is_banned = false`
		case "archived":
			sql += ` AND is_archived = true`
		case "banned":
			sql += ` AND is_banned = true`
		}
	}

	if userTypeFilter != "" {
		sql += ` AND user_type_key = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, userTypeFilter)
		argIndex++
	}

	sql += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	fmt.Println("GetUsers SQL:", sql, "Args:", args)
	err := r.db.DB.Select(&users, sql, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserById(id int) (*UserEntity, error) {
	user := &UserEntity{}
	sql := `SELECT
		*
	FROM users
	WHERE id = $1`
	err := r.db.DB.Get(user, sql, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	sql := `SELECT
		*
	FROM users
	WHERE email = $1`
	err := r.db.DB.Get(user, sql, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateNewUser(dto *models.UserCreateDTO) (*UserEntity, error) {
	sql := `INSERT INTO users (email, display_name, password_hash)
    VALUES ($1, $2, $3)
    RETURNING id;`
	row := r.db.DB.QueryRow(sql, dto.Email, dto.DisplayName, dto.Password)
	var insertedId int
	err := row.Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	user, err := r.GetUserById(insertedId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) MarkUserVerified(userId *int) (bool, error) {
	sql := `UPDATE users SET is_verified = true, modified_at = NOW() WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UpdateUser(dto *models.UserUpdateDTO, userId *int) (*UserEntity, error) {
	sql := `UPDATE users
		SET
			email = $1,
			display_name = $2,
			avatar_url = $3,
			profile_text = $4,
			user_type_key = $5,
			modified_at = NOW()
		WHERE id = $6;`
	_, err := r.db.DB.Exec(sql, dto.Email, dto.DisplayName, dto.AvatarURL, dto.ProfileText, dto.UserTypeKey, &userId)
	if err != nil {
		return nil, err
	}
	user, err := r.GetUserById(*userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUserLastLogin(userId *int) (bool, error) {
	sql := `UPDATE users SET last_login = NOW() WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UpdateUserPassword(userId *int, password string) (bool, error) {
	sql := `UPDATE users SET password_hash = $1 WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, password, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) BanUserByIdWithReason(userId *int, reason string) (bool, error) {
	sql := `UPDATE users
		SET
			is_banned = true,
			ban_reason = $1
		WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, reason, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UnbanUserById(userId *int) (bool, error) {
	sql := `UPDATE users
		SET
			is_banned = false,
			ban_reason = ''
		WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) SetUserTypeKey(userId *int, key string) (bool, error) {
	sql := `UPDATE users
		SET
			user_type_id = $1
		WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, key, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Metrics methods
func (r *UserRepository) GetTotalUsersCount() (*int64, error) {
	count := new(int64)
	sql := `SELECT COUNT(*) as count FROM users`
	err := r.db.DB.Get(count, sql)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (r *UserRepository) GetActiveUsersCount() (*int64, error) {
	count := new(int64)
	sql := `SELECT COUNT(*) as count FROM users WHERE last_login >= NOW() - INTERVAL '24 hours'`
	err := r.db.DB.Get(count, sql)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (r *UserRepository) GetNewUsersCount() (*int64, error) {
	count := new(int64)
	sql := `SELECT COUNT(*) as count FROM users WHERE created_at >= NOW() - INTERVAL '7 days'`
	err := r.db.DB.Get(count, sql)
	if err != nil {
		return nil, err
	}
	return count, nil
}
