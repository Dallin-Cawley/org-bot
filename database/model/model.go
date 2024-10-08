package model

import (
	"orgBot/database/utils"
)

// Model provides an interface for a table in the database. This enforces unique utils.Statement's
// for each required operation on that table.
type Model interface {
	GetInsertStatement() *utils.Statement
	GetDeleteStatement() *utils.Statement
	GetReadStatement() *utils.Statement
}
