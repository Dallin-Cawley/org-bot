package discordapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func FindBotRole(guildID string, session *discordgo.Session) (*discordgo.Role, error) {
	return FindRoleByName("Org Bot", guildID, session)
}

func FindRoleByName(roleName string, guildID string, session *discordgo.Session) (*discordgo.Role, error) {
	roles, err := session.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}

	return nil, fmt.Errorf("no role found with name [ %s ]", roleName)
}
