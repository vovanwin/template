package user

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"template/internal/entity"
	"template/internal/repository/repoerrs"
	"template/pkg/postgres"
	"time"
)

var _ UserRepo = (*PgUserRepo)(nil)

type (
	UserRepo interface {
		GetByID(ctx context.Context, id int64) (user entity.User, err error)
		GetByLogin(ctx context.Context, login string) (user entity.User, err error)
		Delete(ctx context.Context, id int64) (err error)
	}
	PgUserRepo struct {
		*postgres.Postgres
		*slog.Logger
	}
)

func NewPgUserRepo(engine *postgres.Postgres, log *slog.Logger) UserRepo {
	if engine == nil {
		panic("База данных is null")
	}
	return &PgUserRepo{
		Postgres: engine,
		Logger:   log,
	}
}

func (pg *PgUserRepo) GetByID(ctx context.Context, id int64) (user entity.User, err error) {
	sql, args, _ := pg.Builder.
		Select("id, login, tenant, last_login, last_logout, users_status_id, users_role_id, delete, delete_timestamp").
		From("users").
		Where("id = ?", 1).
		Where("delete_timestamp IS NOT NULL").
		ToSql()

	err = pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Tenant,
		&user.LastLogin,
		&user.LastLogout,
		&user.UsersStatusId,
		&user.UsersRoleId,
		&user.Delete,
		&user.DeleteTimestamp,
	)

	//pg.Logger.With("X-Request-ID", 1123).Error("Usage Statistics", sl.Err(err))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetByID - pg.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (pg *PgUserRepo) GetByLogin(ctx context.Context, login string) (user entity.User, err error) {
	sql, args, _ := pg.Builder.
		Select("id, login, tenant, last_login, last_logout, users_status_id, users_role_id, delete, delete_timestamp").
		From("users").
		Where("login = ?", login).
		ToSql()

	err = pg.Pool.QueryRow(ctx, sql, args...).Scan(&user)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetByLogin - pg.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (pg *PgUserRepo) Delete(ctx context.Context, id int64) (err error) {
	sql, args, _ := pg.Builder.
		Update("users").
		SetMap(sq.Eq{
			"delete":           true,
			"delete_timestamp": time.Now(),
		}).
		Where("id = ?", id).
		Where("delete_timestamp not NULL").
		ToSql()

	_, err = pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.Delete - pg.Pool.Exec: %v", err)
	}

	return nil
}
