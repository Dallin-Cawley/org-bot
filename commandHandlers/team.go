package commandHandlers

import (
	"context"
	"fmt"
	"log"

	"orgBot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func HandleTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	action := interaction.ApplicationCommandData().Options[0].StringValue()
	teamName := interaction.ApplicationCommandData().Options[1].StringValue()

	switch action {
	case "onboard":
		onboardTeam(session, interaction, teamName)
	case "offboard":
		deleteTeam(session, interaction, teamName)
	default:
		if err := respondToInteraction(fmt.Sprintf("[ %s ] action is not supported", action), interaction, session); err != nil {
			log.Printf("unable to respond to team interaction due to error: %s\n", err.Error())
		}
	}
}

func onboardTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate, teamName string) {
	role, err := createTeamRole(teamName, session, interaction)
	if err != nil {
		log.Printf("unable to create role with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("Role creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction due to team role creation with error: %s\n", err.Error())
		}

		return
	}

	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		log.Printf("unable to begin transaction with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("unable to begin transaction when onboarding team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting transaction begin failure due to error: %s\n", err.Error())
		}

		return
	}

	botRoleID, err := getBotRoleID(interaction.GuildID, transaction)
	if err != nil {
		log.Printf("unable to get bot role with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("bot role retrieval failed [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting bot role retrieval failure due to error: %s\n", err.Error())
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	if err = transaction.Commit(context.Background()); err != nil {
		log.Printf("unable to commit transaction with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("transaction commit failed while onboarding team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting bot role retrieval failure due to error: %s\n", err.Error())
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	category, err := createTeamCategory(teamName, role, botRoleID, session, interaction)
	if err != nil {
		log.Printf("unable to create category with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("category creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team category creation due to error: %s\n", err.Error())
		}

		return
	}

	if _, err = createChannelInCategoryWithRole(category, "discussion", role, botRoleID, discordgo.ChannelTypeGuildForum, session, interaction); err != nil {
		log.Printf("unable to create forum with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("forum creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team forum creation due to error: %s\n", err.Error())
		}

		return
	}

	if _, err = createChannelInCategoryWithRole(category, "🔉team chat", role, botRoleID, discordgo.ChannelTypeGuildVoice, session, interaction); err != nil {
		log.Printf("unable to create team voice chat with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("team voice chat creation failed for team [ %s ]", teamName), interaction, session); err != nil {
			log.Printf("unable to respond to interaction reporting team voice chat creation due to error: %s\n", err.Error())
		}

		return
	}

	if err = respondToInteraction(fmt.Sprintf("Successfully onboarded team [ %s ]", teamName), interaction, session); err != nil {
		log.Printf("unable to respond to team addJoinToCreateChannel interaction with success due to error: %s\n", err.Error())
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

func createTeamCategory(teamName string, role *discordgo.Role, botRoleID string, session *discordgo.Session, interaction *discordgo.InteractionCreate) (*discordgo.Channel, error) {
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
			{
				ID:    botRoleID,
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

func createChannelInCategoryWithRole(category *discordgo.Channel, name string, role *discordgo.Role, botRoleID string, channelType discordgo.ChannelType, session *discordgo.Session, interaction *discordgo.InteractionCreate) (*discordgo.Channel, error) {
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
			{
				ID:    botRoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
			},
		},
		ParentID: category.ID,
		NSFW:     false,
	}

	createdChannel, err := session.GuildChannelCreateComplex(interaction.GuildID, channelData)
	if err != nil {
		return nil, err
	}

	return createdChannel, nil
}

func createChannelInCategory(categoryID string, name string, channelType discordgo.ChannelType, guildID string, session *discordgo.Session) (*discordgo.Channel, error) {
	channelData := discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     channelType,
		ParentID: categoryID,
		NSFW:     false,
	}

	createdChannel, err := session.GuildChannelCreateComplex(guildID, channelData)
	if err != nil {
		return nil, err
	}

	return createdChannel, nil
}

func deleteTeam(session *discordgo.Session, interaction *discordgo.InteractionCreate, teamName string) {
	log.Println("deleting category")
	if err := deleteCategory(session, interaction, teamName); err != nil {
		log.Printf("unable to deleteChannel category with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("Unable to get channels for category deletion in guild [ %s ]", interaction.GuildID), interaction, session); err != nil {
			log.Printf("unable to respond to guild channel retrieval action due to error: %s\n", err.Error())
		}

		return
	}

	log.Println("deleting role")
	if err := deleteRole(session, interaction.GuildID, teamName); err != nil {
		log.Printf("unable to deleteChannel role with error: %s\n", err.Error())
		if err = respondToInteraction(fmt.Sprintf("Unable to deleteChannel role in guild [ %s ]", interaction.GuildID), interaction, session); err != nil {
			log.Printf("unable to respond to role deletion due to error: %s\n", err.Error())
		}

		return
	}

	if err := respondToInteraction(fmt.Sprintf("Successfully deleted team [ %s ]", teamName), interaction, session); err != nil {
		log.Printf("unable to respond to team deleteChannel action due to error: %s\n", err.Error())
		return
	}

	log.Printf("team [ %s ] deleted successfully\n", teamName)
}

func deleteRole(session *discordgo.Session, guildID string, teamName string) error {
	roles, err := session.GuildRoles(guildID)
	if err != nil {
		return err
	}

	var roleFound bool
	for _, role := range roles {
		if roleFound = caseInsensitiveEquality(role.Name, teamName); roleFound {
			if err = session.GuildRoleDelete(guildID, role.ID); err != nil {
				return err
			}
		}
	}

	if !roleFound {
		return fmt.Errorf("role not found for team [ %s ]", teamName)
	}

	return nil
}

func deleteCategory(session *discordgo.Session, interaction *discordgo.InteractionCreate, categoryName string) error {
	log.Println("getting channels")
	channels, err := session.GuildChannels(interaction.GuildID)
	if err != nil {
		return err
	}

	log.Printf("searching for category [ %s ]\n", categoryName)
	categoryFound := false
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if categoryFound = caseInsensitiveEquality(channel.Name, categoryName); categoryFound {
			log.Println("deleting channels in category")
			if err = deleteChannelsInCategory(channels, channel.ID, session); err != nil {
				return err
			}
		}
	}

	if !categoryFound {
		log.Printf("category [ %s ] not found for deletion\n", categoryName)
		return fmt.Errorf("no category found with name [ %s ]", categoryName)
	}

	log.Printf("successfully deleted category [ %s ]\n", categoryName)
	return nil
}

func deleteChannelsInCategory(discordChannels []*discordgo.Channel, categoryID string, session *discordgo.Session) error {
	var channels []chan *ChannelResponse

	for _, channel := range discordChannels {
		if channel.ParentID == categoryID {
			theChannel := make(chan *ChannelResponse)
			go deleteChannel(theChannel, channel.ID, session)

			channels = append(channels, theChannel)
		}
	}

	var errorMessage string
	for i, channel := range channels {
		log.Printf("awaiting response from channel [ %d ]\n", i)
		response := <-channel

		if !response.Success {
			errorMessage = fmt.Sprintf("%s\n%s", errorMessage, response.Message)
		}
	}

	if errorMessage != "" {
		log.Printf("There was an error deleting channels: %s\n", errorMessage)
		return errors.New(errorMessage)
	}

	log.Println("deleting category")
	channel := make(chan *ChannelResponse)
	go deleteChannel(channel, categoryID, session)

	log.Println("awaiting response from category deletion")
	response := <-channel
	if !response.Success {
		log.Println("category deletion failed with error: " + response.Message)
		return errors.New(response.Message)
	}

	log.Println("successfully deleted channels in category")
	return nil
}

func deleteChannel(channel chan *ChannelResponse, channelID string, session *discordgo.Session) {
	log.Printf("deleting channel [ %s ]\n", channelID)
	if _, err := session.ChannelDelete(channelID); err != nil {
		log.Printf("unable to deleteChannel channel [ %s ] with error: %s\n", channelID, err.Error())
		channel <- &ChannelResponse{Success: false, Message: err.Error()}
		return
	}

	log.Println("successfully deleted channel")
	channel <- &ChannelResponse{Success: true, Message: ""}
}

type ChannelResponse struct {
	Success bool
	Message string
}
