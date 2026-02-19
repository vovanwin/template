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

type OAuthAccount struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Provider   string
	ProviderID string
	Email      string
	CreatedAt  time.Time
}

type OAuthRepo struct {
	pg *postgres.Postgres
}

func NewOAuthRepo(pg *postgres.Postgres) *OAuthRepo {
	return &OAuthRepo{pg: pg}
}

func (r *OAuthRepo) Create(ctx context.Context, userID uuid.UUID, provider, providerID, email string) (*OAuthAccount, error) {
	query, args, err := r.pg.Builder.
		Insert("oauth_accounts").
		Columns("user_id", "provider", "provider_id", "email").
		Values(userID, provider, providerID, email).
		Suffix("RETURNING id, user_id, provider, provider_id, COALESCE(email, ''), created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var a OAuthAccount
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.UserID, &a.Provider, &a.ProviderID, &a.Email, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert oauth account: %w", err)
	}
	return &a, nil
}

func (r *OAuthRepo) GetByProvider(ctx context.Context, provider, providerID string) (*OAuthAccount, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "provider", "provider_id", "COALESCE(email, '')", "created_at").
		From("oauth_accounts").
		Where(squirrel.Eq{"provider": provider, "provider_id": providerID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var a OAuthAccount
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.UserID, &a.Provider, &a.ProviderID, &a.Email, &a.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get oauth account: %w", err)
	}
	return &a, nil
}

func (r *OAuthRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]OAuthAccount, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "provider", "provider_id", "COALESCE(email, '')", "created_at").
		From("oauth_accounts").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list oauth accounts: %w", err)
	}
	defer rows.Close()

	var accounts []OAuthAccount
	for rows.Next() {
		var a OAuthAccount
		if err := rows.Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderID, &a.Email, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan oauth account: %w", err)
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *OAuthRepo) DeleteByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) error {
	query, args, err := r.pg.Builder.
		Delete("oauth_accounts").
		Where(squirrel.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete oauth account: %w", err)
	}
	return nil
}
