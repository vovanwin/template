package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vovanwin/template/app/dbsqlc"
	"github.com/vovanwin/template/app/pkg/utils/hasher"
)

// User представляет пользователя для бизнес-логики
type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	FirstName     *string   `json:"first_name,omitempty"`
	LastName      *string   `json:"last_name,omitempty"`
	Role          string    `json:"role"`
	TenantID      *string   `json:"tenant_id,omitempty"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	Settings      string    `json:"settings,omitempty"`
	Components    string    `json:"components,omitempty"`
}

type UsersRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	ValidatePassword(ctx context.Context, email, password string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type PostgresUsersRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewPostgresUsersRepository(db *pgxpool.Pool) UsersRepository {
	return &PostgresUsersRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *PostgresUsersRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.dbUserToUser(dbUser), nil
}

func (r *PostgresUsersRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	pgUUID := pgtype.UUID{}
	err := pgUUID.Scan(id)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.dbUserToUser(dbUser), nil
}

func (r *PostgresUsersRepository) ValidatePassword(ctx context.Context, email, password string) (*User, error) {
	user, err := r.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // пользователь не найден
	}

	// Проверяем пароль через существующий hasher
	valid, err := hasher.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, nil // неверный пароль
	}

	return user, nil
}

func (r *PostgresUsersRepository) Create(ctx context.Context, user *User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	pgUUID := pgtype.UUID{}
	err := pgUUID.Scan(user.ID)
	if err != nil {
		return err
	}

	params := dbsqlc.CreateUserParams{
		ID:            pgUUID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          &user.Role,
		TenantID:      user.TenantID,
		IsActive:      &user.IsActive,
		EmailVerified: &user.EmailVerified,
		Settings:      []byte(user.Settings),
		Components:    []byte(user.Components),
	}

	_, err = r.queries.CreateUser(ctx, params)
	return err
}

func (r *PostgresUsersRepository) Update(ctx context.Context, user *User) error {
	pgUUID := pgtype.UUID{}
	err := pgUUID.Scan(user.ID)
	if err != nil {
		return err
	}

	params := dbsqlc.UpdateUserParams{
		ID:            pgUUID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          &user.Role,
		TenantID:      user.TenantID,
		IsActive:      &user.IsActive,
		EmailVerified: &user.EmailVerified,
		Settings:      []byte(user.Settings),
		Components:    []byte(user.Components),
	}

	_, err = r.queries.UpdateUser(ctx, params)
	return err
}

// dbUserToUser конвертирует sqlc User в доменную модель User
func (r *PostgresUsersRepository) dbUserToUser(dbUser *dbsqlc.Users) *User {
	user := &User{
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FirstName:    dbUser.FirstName,
		LastName:     dbUser.LastName,
		TenantID:     dbUser.TenantID,
		Settings:     string(dbUser.Settings),
		Components:   string(dbUser.Components),
	}

	// Конвертируем UUID
	if dbUser.ID.Valid {
		copy(user.ID[:], dbUser.ID.Bytes[:])
	}

	// Конвертируем поля с default значениями
	if dbUser.Role != nil {
		user.Role = *dbUser.Role
	} else {
		user.Role = "user"
	}

	if dbUser.IsActive != nil {
		user.IsActive = *dbUser.IsActive
	} else {
		user.IsActive = true
	}

	if dbUser.EmailVerified != nil {
		user.EmailVerified = *dbUser.EmailVerified
	} else {
		user.EmailVerified = false
	}

	return user
}

// HashPassword хеширует пароль с помощью существующего hasher
func HashPassword(password string) (string, error) {
	return hasher.CreateHash(password, hasher.DefaultParams)
}
