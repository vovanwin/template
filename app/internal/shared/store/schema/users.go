package schema

import (
	"github.com/vovanwin/template/internal/shared/types"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Problem holds the schema definition for the Problem entity.
type User struct {
	ent.Schema
}

// Fields of the Problem.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.UserID{}).Default(types.NewUserID).Unique().Immutable(),
		field.String("email").Unique(),
		field.Time("deleted_at").Optional(),
		field.Time("updated_at").Optional().Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Problem.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
