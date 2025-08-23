package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 1
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

	Builder   squirrel.StatementBuilderType
	Pool      PgxPool
	TxManager *TxManager
}

//go:generate options-gen -out-filename=pgx_options.gen.go -from-struct=Options
type Options struct {
	host         string `option:"mandatory" validate:"required"`
	user         string `option:"mandatory" validate:"required"`
	password     string `option:"mandatory" validate:"required"`
	db           string `option:"mandatory" validate:"required"`
	port         string `option:"mandatory" validate:"required"`
	scheme       string `option:"mandatory" validate:"required"`
	isProduction bool   `option:"mandatory"`
}

func New(opts Options) (*Postgres, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options error: %w", err)
	}
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?search_path=%s",
		opts.user,
		opts.password,
		opts.host,
		opts.port,
		opts.db,
		opts.scheme,
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

	// Включает логирование запросов в БД при дебаг режиме
	if opts.isProduction {
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}
	if !opts.isProduction {
		tracer := &tracelog.TraceLog{
			LogLevel: tracelog.LogLevelTrace,
		}
		poolConfig.ConnConfig.Tracer = tracer
		//poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

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

	conn, err := pg.Pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pgdb - New - pgxpool.Acquire: %w", err)
	}
	pg.TxManager = NewTxManager(conn.Conn())
	conn.Release()

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// RunInTx runs the provided function within a transaction.
func (p *Postgres) RunInTx(ctx context.Context, f func(context.Context) error) error {
	return p.TxManager.RunInTx(ctx, f)
}
