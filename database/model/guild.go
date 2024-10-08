package model

import (
	"time"

	"orgBot/database/utils"

	"github.com/jackc/pgx/v5"
)

type Guild struct {
	GuildID   string    `db:"guild_id"`
	BotRoleID string    `db:"bot_role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// MakeGuild creates a new model.Guild with the provided guildID and botRoleID
func MakeGuild(guildID string, botRoleID string) Guild {
	return Guild{GuildID: guildID, BotRoleID: botRoleID}
}

// GetInsertStatement implements the model.Model interface to create an appropriate utils.Statement for a new
// model.Guild record in the database.
func (guild Guild) GetInsertStatement() *utils.Statement {
	argsList := []string{"guild_id", "bot_role_id"}
	namedArgs := pgx.NamedArgs{
		"guild_id":    guild.GuildID,
		"bot_role_id": guild.BotRoleID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "new_guild", argsList, namedArgs)
}

// GetReadStatement implements the model.Model interface to create an appropriate utils.Statement for reading
// a model.Guild record from the database.
func (guild Guild) GetReadStatement() *utils.Statement {
	argsList := []string{"guild_id"}
	namedArgs := pgx.NamedArgs{
		"guild_id": guild.GuildID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "read_guild", argsList, namedArgs)
}

// GetDeleteStatement implements the model.Model interface to create an appropriate utils.Statement for
// deleting a model.Guild record from the database.
func (guild Guild) GetDeleteStatement() *utils.Statement {
	argsList := []string{"guild_id"}
	namedArgs := pgx.NamedArgs{
		"guild_id": guild.GuildID,
	}

	return utils.NewMethodStatement("SELECT * FROM", "delete_guild", argsList, namedArgs)
}
