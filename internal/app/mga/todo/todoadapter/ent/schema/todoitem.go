package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TodoItem holds the schema definition for the TodoItem entity.
type TodoItem struct {
	ent.Schema
}

// Fields of the TodoItem.
func (TodoItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid").
			MaxLen(26).
			NotEmpty().
			Unique().
			Immutable(),
		field.Text("title"),
		field.Bool("completed"),
		field.Int("order"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the TodoItem.
func (TodoItem) Edges() []ent.Edge {
	return nil
}
