// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeSession holds the string denoting the session edge name in mutations.
	EdgeSession = "session"
	// SessionFieldID holds the string denoting the ID field of the Session.
	SessionFieldID = "token"
	// Table holds the table name of the user in the database.
	Table = "users"
	// SessionTable is the table that holds the session relation/edge.
	SessionTable = "sessions"
	// SessionInverseTable is the table name for the Session entity.
	// It exists in this package in order to avoid circular dependency with the "session" package.
	SessionInverseTable = "sessions"
	// SessionColumn is the table column denoting the session relation/edge.
	SessionColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// BySessionField orders the results by session field.
func BySessionField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSessionStep(), sql.OrderByField(field, opts...))
	}
}
func newSessionStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SessionInverseTable, SessionFieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, SessionTable, SessionColumn),
	)
}
