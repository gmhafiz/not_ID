package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the User.
func (Session) Fields() []ent.Field {
	return []ent.Field{

		//field.ID("token"),
		field.String("id").StorageKey("token"),
		field.Uint64("user_id").Optional().Nillable(),
		field.Time("created_at"),
	}
}

// Edges of the User.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("session").
			Unique(),
	}
}
