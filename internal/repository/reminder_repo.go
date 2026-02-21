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

type Reminder struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	RemindAt    time.Time
	WorkflowID  string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReminderRepo struct {
	pg *postgres.Postgres
}

func NewReminderRepo(pg *postgres.Postgres) *ReminderRepo {
	return &ReminderRepo{pg: pg}
}

func (r *ReminderRepo) Create(ctx context.Context, userID uuid.UUID, title, description string, remindAt time.Time) (*Reminder, error) {
	query, args, err := r.pg.Builder.
		Insert("reminders").
		Columns("user_id", "title", "description", "remind_at").
		Values(userID, title, description, remindAt).
		Suffix("RETURNING id, user_id, title, description, remind_at, COALESCE(workflow_id, ''), status, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var rem Reminder
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&rem.ID, &rem.UserID, &rem.Title, &rem.Description, &rem.RemindAt,
		&rem.WorkflowID, &rem.Status, &rem.CreatedAt, &rem.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert reminder: %w", err)
	}
	return &rem, nil
}

func (r *ReminderRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]Reminder, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "title", "description", "remind_at", "COALESCE(workflow_id, '')", "status", "created_at", "updated_at").
		From("reminders").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("remind_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query reminders: %w", err)
	}
	defer rows.Close()

	var result []Reminder
	for rows.Next() {
		var rem Reminder
		if err := rows.Scan(
			&rem.ID, &rem.UserID, &rem.Title, &rem.Description, &rem.RemindAt,
			&rem.WorkflowID, &rem.Status, &rem.CreatedAt, &rem.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		result = append(result, rem)
	}
	return result, nil
}

func (r *ReminderRepo) GetByID(ctx context.Context, id uuid.UUID) (*Reminder, error) {
	query, args, err := r.pg.Builder.
		Select("id", "user_id", "title", "description", "remind_at", "COALESCE(workflow_id, '')", "status", "created_at", "updated_at").
		From("reminders").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var rem Reminder
	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(
		&rem.ID, &rem.UserID, &rem.Title, &rem.Description, &rem.RemindAt,
		&rem.WorkflowID, &rem.Status, &rem.CreatedAt, &rem.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get reminder: %w", err)
	}
	return &rem, nil
}

func (r *ReminderRepo) UpdateWorkflowID(ctx context.Context, id uuid.UUID, workflowID string) error {
	query, args, err := r.pg.Builder.
		Update("reminders").
		Set("workflow_id", workflowID).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update workflow_id: %w", err)
	}
	return nil
}

func (r *ReminderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query, args, err := r.pg.Builder.
		Update("reminders").
		Set("status", status).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

func (r *ReminderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.pg.Builder.
		Delete("reminders").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete reminder: %w", err)
	}
	return nil
}
