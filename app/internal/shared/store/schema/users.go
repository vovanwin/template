package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/vovanwin/template/internal/shared/types"
	"time"
)

// Users holds the schema definition for the Users entity.
type Users struct {
	ent.Schema
}

// Fields of the Users.
func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.UserID{}).Default(types.NewUserID).Unique().Immutable(),
		field.String("login"),
		field.String("password").MaxLen(255),

		field.Time("deleted_at").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("user", User.Type).Ref("posts").Field("user_id").Unique().Required(),
	}
}
