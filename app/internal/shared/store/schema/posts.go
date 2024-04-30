package schema

import (
	"entgo.io/ent/schema/edge"
	"github.com/vovanwin/template/internal/shared/types"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Problem holds the schema definition for the Problem entity.
type Post struct {
	ent.Schema
}

// Fields of the Problem.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.UserID{}).Default(types.NewUserID).Unique().Immutable(),
		field.String("test"),
		field.String("title"),
		field.String("title1"),
		field.String("title3"),
		field.UUID("user_id", types.UserID{}),
		field.Time("deleted_at").Optional(),
		field.Time("updated_at").Optional().Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Problem.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("posts").Field("user_id").Unique().Required(),
	}
}
