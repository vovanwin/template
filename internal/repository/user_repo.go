package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vovanwin/template/internal/pkg/storage/postgres"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	AvatarURL    string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepo struct {
	pg *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg: pg}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, name string) (*User, error) {
	query, args, err := r.pg.Builder.
		Insert("users").
		Columns("email", "password_hash", "name").
		Values(email, passwordHash, name).
		Suffix("RETURNING id, email, password_hash, name, COALESCE(avatar_url, ''), role, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &u, nil
}

// CreateOAuth создаёт пользователя через OAuth (без пароля).
func (r *UserRepo) CreateOAuth(ctx context.Context, email, name string) (*User, error) {
	query, args, err := r.pg.Builder.
		Insert("users").
		Columns("email", "password_hash", "name").
		Values(email, "", name).
		Suffix("RETURNING id, email, password_hash, COALESCE(name, ''), COALESCE(avatar_url, ''), role, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert oauth user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	query, args, err := r.pg.Builder.
		Select("id", "email", "password_hash", "COALESCE(name, '')", "COALESCE(avatar_url, '')", "role", "is_active", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query, args, err := r.pg.Builder.
		Select("id", "email", "password_hash", "COALESCE(name, '')", "COALESCE(avatar_url, '')", "role", "is_active", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) SetRole(ctx context.Context, id uuid.UUID, role string) (*User, error) {
	query, args, err := r.pg.Builder.
		Update("users").
		Set("role", role).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, email, password_hash, COALESCE(name, ''), COALESCE(avatar_url, ''), role, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("set role: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (*User, error) {
	fields["updated_at"] = squirrel.Expr("NOW()")

	query, args, err := r.pg.Builder.
		Update("users").
		SetMap(fields).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, email, password_hash, COALESCE(name, ''), COALESCE(avatar_url, ''), role, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var u User
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}
