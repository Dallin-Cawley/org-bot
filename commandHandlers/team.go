package commandHandlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	action := interaction.ApplicationCommandData().Options[0].StringValue()
	teamName := interaction.ApplicationCommandData().Options[1].StringValue()

	switch action {
	case "onboard":
		onboardTeam(session, interaction, teamName)
	case "delete":
		deleteTeam(session, interaction, teamName)
	default:
		if err := respondToInteraction(fmt.Sprintf("[ %s ] action is not supported", action), interaction, session); err != nil {
			log.Printf("unable to respond to team interaction due to error: %s", err.Error())
		}
	}
}

func onboardTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate, teamName string) {
	role, err := createTeamRole(teamName, session, interaction)
	if err != nil {
		log.Printf("unable to create role with error: %s", err.Error())
		if err = respondToInteraction(fmt.Sprintf("Role creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction due to team role creation with error: %s", err.Error())
		}

		return
	}

	category, err := createTeamCategory(teamName, role, session, interaction)
	if err != nil {
		log.Printf("unable to create category with error: %s", err.Error())
		if err = respondToInteraction(fmt.Sprintf("category creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team category creation due to error: %s", err.Error())
		}

		return
	}

	if err = createChannelInCategory(category, "discussion", role, discordgo.ChannelTypeGuildForum, session, interaction); err != nil {
		log.Printf("unable to create forum with error: %s", err.Error())
		if err = respondToInteraction(fmt.Sprintf("forum creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team forum creation due to error: %s", err.Error())
		}

		return
	}

	if err = createChannelInCategory(category, "🔉team chat", role, discordgo.ChannelTypeGuildVoice, session, interaction); err != nil {
		log.Printf("unable to create team voice chat with error: %s", err.Error())
		if err = respondToInteraction(fmt.Sprintf("team voice chat creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team voice chat creation due to error: %s", err.Error())
		}

		return
	}

	if err = respondToInteraction(fmt.Sprintf("Successfully onboarded team [ %s ]", teamName), interaction, session); err != nil {
		log.Printf("unable to respond to team add interaction with success due to error: %s", err.Error())
	}
}

func createTeamRole(teamName string, session *discordgo.Session, interaction *discordgo.InteractionCreate) (*discordgo.Role, error) {
	roleParams := discordgo.RoleParams{
		Name: teamName,
	}

	role, err := session.GuildRoleCreate(interaction.GuildID, &roleParams)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func createTeamCategory(teamName string, role *discordgo.Role, session *discordgo.Session, interaction *discordgo.InteractionCreate) (*discordgo.Channel, error) {
	channelData := discordgo.GuildChannelCreateData{
		Name: teamName,
		Type: discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   interaction.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:    role.ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
			},
		},
		NSFW: false,
	}

	category, err := session.GuildChannelCreateComplex(interaction.GuildID, channelData)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func createChannelInCategory(category *discordgo.Channel, name string, role *discordgo.Role, channelType discordgo.ChannelType, session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	channelData := discordgo.GuildChannelCreateData{
		Name: name,
		Type: channelType,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   interaction.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:    role.ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
			},
		},
		ParentID: category.ID,
		NSFW:     false,
	}

	_, err := session.GuildChannelCreateComplex(interaction.GuildID, channelData)
	if err != nil {
		return err
	}

	return nil
}

func deleteTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate, teamName string) {

}

func deleteCategory(session *discordgo.Session, interaction *discordgo.InteractionCreate, categoryName string) {
	channels, err := session.GuildChannels(interaction.GuildID)
	if err != nil {
		if err = respondToInteraction(fmt.Sprintf("Unable to get channels for category deletion in guild [ %s ]", interaction.GuildID), interaction, session); err != nil {
			log.Printf("unable to respond to guild channel retrieval action due to error: %s", err.Error())
		}

		return
	}

	channelFound := false
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if channelFound = channel.Name == categoryName; channelFound {
			if err = deleteChannelsInCategory(channels, )
		}
	}
}

func deleteChannelsInCategory(channels []*discordgo.Channel, categoryID string, session *discordgo.Session) error {
	for _, channel := range channels {
		if channel.ParentID == categoryID {

		}
	}
}
