package model

import (
	"github.com/jackc/pgx/v5"
	"orgBot/database/utils"
	"time"
)

type JoinVCChild struct {
	JoinVCChildID string    `db:"join_vc_child_id"`
	JoinVCID      string    `json:"join_vc_id" db:"join_vc_id"`
	GuildID       string    `json:"guild_id" db:"guild_id"`
	CategoryID    string    `json:"category_id" db:"category_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func MakeJoinVCChild(joinVCChildID string, joinVCID string, guildID string, categoryID string) JoinVCChild {
	return JoinVCChild{JoinVCChildID: joinVCChildID, JoinVCID: joinVCID, GuildID: guildID, CategoryID: categoryID}
}

// GetInsertStatement implements the model.Model interface to create an appropriate utils.Statement for a new
// model.JoinVCChild record in the database.
func (joinVCChild JoinVCChild) GetInsertStatement() *utils.Statement {
	argsList := []string{"join_vc_child_id", "guild_id", "category_id", "join_vc_id"}
	namedArgs := pgx.NamedArgs{
		"join_vc_child_id": joinVCChild.JoinVCChildID,
		"join_vc_id":       joinVCChild.JoinVCID,
		"guild_id":         joinVCChild.GuildID,
		"category_id":      joinVCChild.CategoryID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "new_join_vc_child", argsList, namedArgs)
}

// GetReadStatement implements the model.Model interface to create an appropriate utils.Statement for reading
// a model.JoinVCChild record from the database.
func (joinVCChild JoinVCChild) GetReadStatement() *utils.Statement {
	argsList := []string{"join_vc_child_id"}
	namedArgs := pgx.NamedArgs{
		"join_vc_child_id": joinVCChild.JoinVCChildID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "read_join_vc_child", argsList, namedArgs)
}

// GetDeleteStatement implements the model.Model interface to create an appropriate utils.Statement for
// deleting a model.JoinVCChild record from the database.
func (joinVCChild JoinVCChild) GetDeleteStatement() *utils.Statement {
	argsList := []string{"join_vc_child_id"}
	namedArgs := pgx.NamedArgs{
		"join_vc_child_id": joinVCChild.JoinVCChildID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "delete_join_vc_child", argsList, namedArgs)
}
