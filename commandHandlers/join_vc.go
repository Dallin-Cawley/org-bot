package commandHandlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"orgBot/database"
	"orgBot/database/model"
	databaseUtils "orgBot/database/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	JOIN_VC_CHANNEL_NAME = "🔉Join to Create"
)

func VoiceChannelStatusUpdate(session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) {
	if interaction.BeforeUpdate == nil || len(interaction.BeforeUpdate.ChannelID) == 0 {
		if len(interaction.ChannelID) > 0 {
			if err := voiceChannelJoin(session, interaction); err != nil {
				log.Printf("unable to process voice channel join with error [ %s ]\n", err.Error())
			}
		}
	} else {
		beforeID := interaction.BeforeUpdate.ChannelID
		afterID := interaction.ChannelID

		if len(beforeID) != 0 && len(afterID) != 0 {
			if err := voiceChannelMove(session, interaction); err != nil {
				log.Printf("unable to process voice channel move with error [ %s ]\n", err.Error())
			}
		} else if len(beforeID) != 0 && len(afterID) == 0 {
			if err := voiceChannelLeave(session, interaction); err != nil {
				log.Printf("unable to process voice channel leave with error [ %s ]\n", err.Error())
			}
		}
	}
}

func voiceChannelJoin(session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		return err
	}

	if err = handleMoveIntoJoinVCChannel(transaction, session, interaction); err != nil {
		return err
	}

	if err = transaction.Commit(context.Background()); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	return nil
}

// voiceChannelMove handles a user moving from one voice channel to another. If the user has moved into a model.JoinVC
// channel, a new channel is created and the user is moved into it. If the user has moved out of a model.JoinVCChild
// channel, the channel is deleted if no one else is present in it.
func voiceChannelMove(session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		return err
	}

	if err = handleMoveIntoJoinVCChannel(transaction, session, interaction); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	if err = handleMoveOutOfJoinVCChildChannel(transaction, session, interaction); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	if err = transaction.Commit(context.Background()); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	return nil
}

// voiceChannelLeave handles a user leaving a voice channel. If the user has left a model.JoinVCChild channel,
// the model.JoinVCChild channel will be deleted if there are no other members in it.
func voiceChannelLeave(session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		return err
	}

	if err = handleMoveOutOfJoinVCChildChannel(transaction, session, interaction); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	if err = transaction.Commit(context.Background()); err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	}

	return nil
}

// handleMoveIntoJoinVCChannel moves the user into a new model.JoinVCChild channel if the joined channel was truly a
// model.JoinVC channel. If not, nothing happens.
func handleMoveIntoJoinVCChannel(transaction pgx.Tx, session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	movedIn, joinVC, err := hasMovedIntoJoinVCChannel(transaction, interaction)
	if err != nil {
		_ = transaction.Rollback(context.Background())
		return err
	} else if movedIn {
		if err = movedIntoJoinVCChannel(joinVC, transaction, session, interaction); err != nil {
			_ = transaction.Rollback(context.Background())
			return err
		}
	}

	return nil
}

// handleMoveOutOfJoinVCChildChannel deletes the model.JoinVCChild channel from the server if no one is left in it.
// If the channel that was moved out of not a model.JoinVCChild, nothing happens.
func handleMoveOutOfJoinVCChildChannel(transaction pgx.Tx, session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	movedOut, joinVCChild, err := hasMovedOutOfJoinVCChildChannel(transaction, interaction)
	if err != nil {
		return err
	} else if movedOut {
		if err = movedOutOfJoinVCChildChannel(joinVCChild, transaction, session); err != nil {
			return err
		}
	}

	return nil
}

// hasMovedIntoJoinVCChannel determines if the joined channel was truly a model.JoinVC channel.
func hasMovedIntoJoinVCChannel(transaction pgx.Tx, interaction *discordgo.VoiceStateUpdate) (bool, model.JoinVC, error) {
	joinVCs, err := database.ReadModel[model.JoinVC](model.MakeJoinVC(interaction.ChannelID, "", ""), transaction)
	if err != nil {
		if databaseUtils.NoRowsRead(err) {
			log.Printf("Joined channel [ %s ] is not a join to create channel\n", interaction.ChannelID)
			return false, model.JoinVC{}, nil
		}

		return false, model.JoinVC{}, err
	}

	return true, joinVCs[0], nil
}

// hasMovedOutOfJoinVCChildChannel determines if the left channel was truly a model.JoinVCChild channel.
func hasMovedOutOfJoinVCChildChannel(transaction pgx.Tx, interaction *discordgo.VoiceStateUpdate) (bool, model.JoinVCChild, error) {
	joinVCChildren, err := database.ReadModel[model.JoinVCChild](model.MakeJoinVCChild(interaction.BeforeUpdate.ChannelID, model.JoinVC{}), transaction)
	if err != nil {
		if databaseUtils.NoRowsRead(err) {
			log.Printf("Joined channel [ %s ] is not a join to create child channel\n", interaction.ChannelID)
			return false, model.JoinVCChild{}, nil
		}

		return false, model.JoinVCChild{}, err
	}

	return true, joinVCChildren[0], nil
}

// movedOutOfJoinVCChildChannel deletes the model.JoinVCChild channel if no users are left inside.
func movedOutOfJoinVCChildChannel(joinVCChild model.JoinVCChild, transaction pgx.Tx, session *discordgo.Session) error {
	numMembers, err := numMembersInVoiceChannel(joinVCChild.JoinVCChildID, joinVCChild.GuildID, session)
	if err != nil {
		return err
	}

	if numMembers == 0 {
		return deleteJoinVCChildChannel(joinVCChild, transaction, session)
	}

	return nil
}

// deleteJoinVCChildChannel deletes the model.JoinVCChild channel from both the server and the database.
func deleteJoinVCChildChannel(joinVCChild model.JoinVCChild, transaction pgx.Tx, session *discordgo.Session) error {
	if _, err := session.ChannelDelete(joinVCChild.JoinVCChildID); err != nil {
		return err
	}

	if _, err := database.DeleteModel[model.JoinVCChild](joinVCChild, transaction); err != nil {
		return err
	}

	return nil
}

// movedIntoJoinVCChannel creates a new model.JoinVCChild channel and moves the user into it.
func movedIntoJoinVCChannel(joinVC model.JoinVC, transaction pgx.Tx, session *discordgo.Session, interaction *discordgo.VoiceStateUpdate) error {
	channelName := fmt.Sprintf("%s's channel", interaction.Member.DisplayName())

	createdChannel, err := createChannelInCategory(joinVC.CategoryID, channelName, discordgo.ChannelTypeGuildVoice, interaction.GuildID, session)
	if err != nil {
		return err
	}

	if _, err = database.InsertModel[model.JoinVCChild](model.MakeJoinVCChild(createdChannel.ID, joinVC), transaction); err != nil {
		return err
	}

	if err = session.GuildMemberMove(interaction.GuildID, interaction.UserID, &createdChannel.ID); err != nil {
		return err
	}

	return nil
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

	if channels, err := database.ReadModel[model.JoinVC](model.MakeJoinVC("", "", category.ID), transaction); err != nil {
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

	createdChannel, err := createChannelInCategory(category.ID, JOIN_VC_CHANNEL_NAME, discordgo.ChannelTypeGuildVoice, interaction.GuildID, session)
	if err != nil {
		log.Printf("unable to create join to create voice channel in category [ %s ]\n", err)
		if err = respondToInteraction(fmt.Sprintf("unable to create join to create voice channel in category [ %s ]", category.Name), interaction, session); err != nil {
			log.Printf("unable to respond when creating join to create voice channel with error [ %s ]\n", err)
		}

		return
	}

	joinVC := model.MakeJoinVC(createdChannel.ID, interaction.GuildID, category.ID)
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

	if _, err = database.DeleteModel[model.JoinVC](model.MakeJoinVC(channel.ID, "", ""), transaction); err != nil {
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
