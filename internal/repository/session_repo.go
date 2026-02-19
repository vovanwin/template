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

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	IP               string
	UserAgent        string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type SessionRepo struct {
	pg *postgres.Postgres
}

func NewSessionRepo(pg *postgres.Postgres) *SessionRepo {
	return &SessionRepo{pg: pg}
}

func (r *SessionRepo) Create(ctx context.Context, userID uuid.UUID, refreshTokenHash, ip, userAgent string, expiresAt time.Time) (*Session, error) {
	query, args, err := r.pg.Builder.
		Insert("sessions").
		Columns("user_id", "refresh_token_hash", "ip", "user_agent", "expires_at").
		Values(userID, refreshTokenHash, ip, userAgent, expiresAt).
		Suffix("RETURNING id, user_id, refresh_token_hash, ip, user_agent, expires_at, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var s Session
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &s.IP, &s.UserAgent, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}
	return &s, nil
}

func (r *SessionRepo) GetByTokenHash(ctx context.Context, hash string) (*Session, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "refresh_token_hash", "ip", "user_agent", "expires_at", "created_at").
		From("sessions").
		Where(squirrel.Eq{"refresh_token_hash": hash}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var s Session
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &s.IP, &s.UserAgent, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get session by token hash: %w", err)
	}
	return &s, nil
}

func (r *SessionRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "refresh_token_hash", "ip", "user_agent", "expires_at", "created_at").
		From("sessions").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.RefreshTokenHash, &s.IP, &s.UserAgent, &s.ExpiresAt, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionRepo) Delete(ctx context.Context, sessionID uuid.UUID) error {
	query, args, err := r.pg.Builder.
		Delete("sessions").
		Where(squirrel.Eq{"id": sessionID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteByTokenHash(ctx context.Context, hash string) error {
	query, args, err := r.pg.Builder.
		Delete("sessions").
		Where(squirrel.Eq{"refresh_token_hash": hash}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete session by token hash: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query, args, err := r.pg.Builder.
		Delete("sessions").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}
	return nil
}
