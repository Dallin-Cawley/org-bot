package model

import (
	"github.com/jackc/pgx/v5"
	"orgBot/database/utils"
	"time"
)

type JoinVC struct {
	JoinVCID   string    `json:"join_vc_id" db:"join_vc_id"`
	GuildID    string    `json:"guild_id" db:"guild_id"`
	CategoryID string    `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// NewJoinVC creates a new model.JoinVC
func NewJoinVC(joinVCID string, guildID string, categoryID string) JoinVC {
	return JoinVC{JoinVCID: joinVCID, GuildID: guildID, CategoryID: categoryID}
}

// GetInsertStatement implements the model.Model interface to create an appropriate utils.Statement for a new
// model.JoinVC record in the database.
func (joinVCChild JoinVC) GetInsertStatement() *utils.Statement {
	argsList := []string{"join_vc_id", "guild_id", "category_id"}
	namedArgs := pgx.NamedArgs{
		"join_vc_id":  joinVCChild.JoinVCID,
		"guild_id":    joinVCChild.GuildID,
		"category_id": joinVCChild.CategoryID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "new_join_vc", argsList, namedArgs)
}

// GetReadStatement implements the model.Model interface to create an appropriate utils.Statement for reading
// a model.JoinVC record from the database.
func (joinVCChild JoinVC) GetReadStatement() *utils.Statement {
	var argsList []string
	var namedArgs pgx.NamedArgs
	var methodName string

	if joinVCChild.JoinVCID != "" {
		argsList = []string{"join_vc_id"}
		namedArgs = pgx.NamedArgs{
			"join_vc_id": joinVCChild.JoinVCID,
		}
		methodName = "read_join_vc"
	} else {
		argsList = []string{"category_id"}
		namedArgs = pgx.NamedArgs{
			"category_id": joinVCChild.CategoryID,
		}
		methodName = "read_join_vc_in_category"
	}

	return utils.NewMethodStatement("SELECT * FROM", methodName, argsList, namedArgs)
}

// GetDeleteStatement implements the model.Model interface to create an appropriate utils.Statement for
// deleting a model.JoinVC record from the database.
func (joinVCChild JoinVC) GetDeleteStatement() *utils.Statement {
	argsList := []string{"join_vc_id"}
	namedArgs := pgx.NamedArgs{
		"join_vc_id": joinVCChild.JoinVCID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "delete_join_vc", argsList, namedArgs)
}
