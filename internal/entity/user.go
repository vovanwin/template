package entity

type User struct {
	ID              int    `db:"id" json:"id"`
	Login           string `db:"login" json:"login"`
	Tenant          string `db:"tenant" json:"tenant"`
	LastLogin       string `db:"last_login" json:"last_login"`
	LastLogout      string `db:"last_logout" json:"last_logout"`
	UsersStatusId   int    `db:"users_status_id" json:"status_id"`
	UsersRoleId     int    `db:"users_role_id" json:"role_id"`
	Delete          bool   `db:"delete" json:"delete"`
	DeleteTimestamp *int   `db:"delete_timestamp" json:"delete_timestamp"`
}
