package commandHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"orgBot/database"
	"orgBot/database/model"
	databaseUtils "orgBot/database/utils"
	"orgBot/discordapi"

	"github.com/bwmarrin/discordgo"
)

var (
	configDir = filepath.Join(os.Getenv("APP_DIR"), "config")
)

func GuildCreate(session *discordgo.Session, interaction *discordgo.GuildCreate) {
	log.Printf("Bot entered guild %s\n", interaction.Name)

	slashCommands, err := initSlashCommands()
	if err != nil {
		log.Printf("Error initializing slash commands: %s\n", err.Error())
		return
	}

	if _, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, interaction.ID, slashCommands); err != nil {
		log.Printf("Error creating slash commands: %s\n", err.Error())
		return
	}

	botRole, err := discordapi.FindBotRole(interaction.ID, session)
	if err != nil {
		log.Printf("Error finding bot role: %s\n", err.Error())
		return
	}

	guild := model.MakeGuild(interaction.ID, botRole.ID)
	transaction, err := database.BeginTransaction(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %s\n", err.Error())
		return
	}

	_, err = database.InsertModel[model.Guild](guild, transaction)
	if err != nil {
		if !databaseUtils.NoRowsRead(err) {
			log.Printf("Error inserting guild: %s\n", err.Error())
			_ = transaction.Rollback(context.Background())
		}

		return
	}

	if err = transaction.Commit(context.Background()); err != nil {
		log.Printf("Error committing guild creation transaction: %s\n", err.Error())
		_ = transaction.Rollback(context.Background())

		return
	}
}

func initSlashCommands() ([]*discordgo.ApplicationCommand, error) {
	fileBytes, err := os.ReadFile(fmt.Sprintf("%s/slashCommands.json", configDir))
	if err != nil {
		return nil, err
	}

	slashCommands := make([]*discordgo.ApplicationCommand, 0)
	if err = json.Unmarshal(fileBytes, &slashCommands); err != nil {
		return nil, err
	}

	return slashCommands, nil
}
