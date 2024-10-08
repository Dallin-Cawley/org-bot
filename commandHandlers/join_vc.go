package commandHandlers

import (
	"context"
	"fmt"
	"log"
	"orgBot/database"
	"orgBot/database/model"
	databaseUtils "orgBot/database/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	JOIN_VC_CHANNEL_NAME = "🔉Join to Create"
)

func JoinToCreateUserJoin(session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) {
	interaction.
}

func JoinToCreateVoiceChat(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	action := interaction.ApplicationCommandData().Options[0].StringValue()
	category := interaction.ApplicationCommandData().Options[1].ChannelValue(session)

	switch action {
	case "add":
		addJoinToCreateChannel(category, session, interaction)
	case "delete":
		deleteJoinToCreateChannel(category, session, interaction)
	default:
		if err := respondToInteraction(fmt.Sprintf("action [ %s ] is not supported", action), interaction, session); err != nil {
			log.Printf("unable to respond to invalid action [ %s ]\n", action)
		}
	}
}

func addJoinToCreateChannel(category *discordgo.Channel, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		log.Printf("unable to begin transaction when persisting join to create voice channel in database [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("unable to begin transaction when persisting join to create voice channel in category [ %s ]", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when beginning transaction to persist join to create voice channel with error [ %s ]\n", err)
		}

		return
	}

	if channels, err := database.ReadModel[model.JoinVC](model.NewJoinVC("", "", category.ID), transaction); err != nil {
		if !databaseUtils.NoRowsRead(err) {
			log.Printf("attempt to determine if a join to create voice channel in category [ %s ] already exists resulted in an error [ %s ]\n", category.Name, err)
			if err = respondToInteraction(fmt.Sprintf("attempt to determine if a join to create voice channel in category [ %s ] already exists resulted in an error", category.Name), interaction, session); err != nil {
				log.Printf("unable to respond when notifying user of join to create existence attempt failed with error [ %s ]\n", err)
			}

			_ = transaction.Rollback(context.Background())
			return
		}
	} else if len(channels) > 0 {
		log.Printf("there is already a join to create voice channel in category [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("there can only be one join to create voice channel in category [ %s ]. "+
			"If you deleted it manually, please run the delete action on this category to reset.", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when notifying user of already existing join to create voice channel with error [ %s ]\n", err)
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	createdChannel, err := createChannelInCategory(category, JOIN_VC_CHANNEL_NAME, discordgo.ChannelTypeGuildVoice, session, interaction)
	if err != nil {
		log.Printf("unable to create join to create voice channel in category [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("unable to create join to create voice channel in category [ %s ]", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when creating join to create voice channel with error [ %s ]\n", err)
		}

		return
	}

	joinVC := model.NewJoinVC(createdChannel.ID, interaction.GuildID, category.ID)
	if _, err = database.InsertModel[model.JoinVC](joinVC, transaction); err != nil {
		log.Printf("unable to persist join to create voice channel in database [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("unable to persist join to create voice channel in category [ %s ]", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when to persisting join to create voice channel with error [ %s ]\n", err)
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	if err = transaction.Commit(context.Background()); err != nil {
		log.Printf("unable to commit transaction of persisted join to create voice channel in database [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("unable to commit join to create voice channel in category [ %s ]", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when to persisting join to create voice channel with error [ %s ]\n", err)
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	if err = respondToInteraction(fmt.Sprintf("successfully created join to create channel in category [ %s ]", category.Name), interaction, session); err != nil {
		log.Printf("unable to respond when reporting join to create channel creation success with error [ %s ]\n", err)
	}
}

func deleteJoinToCreateChannel(category *discordgo.Channel, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	channel, err := getJoinToCreateChannel(category, session, interaction)
	if err != nil {
		log.Printf("unable to find join to create channel with error [ %s ]\n", err)
		if err = respondToInteraction("unable to find join to create channel", interaction, session); err != nil {
			log.Printf("unable to respond while finding join to create channel with error [ %s ]\n", err)
		}

		return
	}

	_, err = session.ChannelDelete(channel.ID)
	if err != nil {
		log.Printf("unable to delete join to create channel with error [ %s ]\n", err)
		if err = respondToInteraction("unable to delete join to create channel", interaction, session); err != nil {
			log.Printf("unable to respond while deleting join to create voice channel with error [ %s ]\n", err)
		}

		return
	}

	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		log.Printf("unable to begin transaction [ %s ]\n", err)
		if err = respondToInteraction("unable to begin transaction", interaction, session); err != nil {
			log.Printf("unable to respond while beginning transaction with error [ %s ]\n", err)
		}

		return
	}

	if _, err = database.DeleteModel[model.JoinVC](model.NewJoinVC(channel.ID, "", ""), transaction); err != nil {
		if !databaseUtils.NoRowsRead(err) {
			log.Printf("unable to delete join to create channel in database with error [ %s ]\n", err)
			if err = respondToInteraction("unable to delete join to create channel in database", interaction, session); err != nil {
				log.Printf("unable to respond while deleting join to create channel in database with error [ %s ]\n", err)
			}

			_ = transaction.Rollback(context.Background())
			return
		}
	}

	if err = transaction.Commit(context.Background()); err != nil {
		log.Printf("unable to commit transaction of database with error [ %s ]\n", err)
		if err = respondToInteraction("unable to commit transaction", interaction, session); err != nil {
			log.Printf("unable to respond while committing transaction with error [ %s ]\n", err)
		}

		_ = transaction.Rollback(context.Background())
		return
	}

	if err = respondToInteraction("successfully deleted join to create channel", interaction, session); err != nil {
		log.Printf("unable to respond when reporting join to create channel deletion success with error [ %s ]\n", err)
	}
}

func getJoinToCreateChannel(category *discordgo.Channel, session *discordgo.Session, interaction *discordgo.InteractionCreate) (*discordgo.Channel, error) {
	channels, err := session.GuildChannels(interaction.GuildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.ParentID != category.ID {
			continue
		}

		if channel.Name == JOIN_VC_CHANNEL_NAME {
			return channel, nil
		}
	}

	return nil, fmt.Errorf("join to create channel not found in category [ %s ]", category.Name)
}
