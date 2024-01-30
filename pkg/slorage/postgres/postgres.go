package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"log/slog"
	"template/config"
	"time"
)

const (
	defaultMaxPoolSize  = 5
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    PgxPool
}

func New(config config.Config, log *slog.Logger) (*Postgres, error) {

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Scheme,
	)

	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	//Включает логирование запросов в БД при дебаг режиме
	if config.LogLevel == "debug" {
		tracer := &tracelog.TraceLog{
			Logger:   NewLoggerTracer(log),
			LogLevel: tracelog.LogLevelTrace,
		}
		poolConfig.ConnConfig.Tracer = tracer
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		slog.Info("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgxpool.ConnectConfig: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
