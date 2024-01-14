package entity

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID        int              `db:"id" json:"id"`
	Login     string           `db:"login" json:"login"`
	Tenant    string           `db:"tenant" json:"tenant"`
	Password  string           `db:"password" json:"-"`
	StatusId  int              `db:"status_id" json:"status_id"`
	RoleId    int              `db:"role_id" json:"role_id"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
}
