package commandHandlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func AddRoleSync(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	newRoleName := interaction.ApplicationCommandData().Options[0].StringValue()
	copyRole := interaction.ApplicationCommandData().Options[1].RoleValue(session, interaction.GuildID)

	if _, err := session.GuildRoleCreate(interaction.GuildID, getSyncRoleParams(newRoleName, copyRole)); err != nil {
		if err = respondToInteraction(fmt.Sprintf("unable to create role [ %s ] with same permissions as role [ %s ] with error %s", newRoleName, copyRole.Name, err.Error()), interaction, session); err != nil {
			log.Printf("unable to respond to successful copy role creation with error: %s\n", err.Error())
		}

		return
	}

	if err := respondToInteraction(fmt.Sprintf("Successfully created role [ %s ] with same permissions as role [ %s ]", newRoleName, copyRole.Name), interaction, session); err != nil {
		log.Printf("unable to respond to successful copy role creation with error: %s\n", err.Error())
	}
}

func getSyncRoleParams(roleName string, role *discordgo.Role) *discordgo.RoleParams {
	return &discordgo.RoleParams{
		Name:         roleName,
		Color:        nil,
		Hoist:        nil,
		Permissions:  &role.Permissions,
		Mentionable:  nil,
		UnicodeEmoji: nil,
		Icon:         nil,
	}
}
