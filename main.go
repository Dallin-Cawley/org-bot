package main

import (
	"encoding/json"
	"fmt"
	"log"
	"orgBot/commandHandlers"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var (
	configDir = filepath.Join(os.Getenv("APP_DIR"), "config")
)

func main() {
	discordSession, err := discordgo.New("Bot " + os.Getenv("DISCORD_API_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	discordSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.GuildCreate) {
		log.Printf("Bot entered guild %s\n", i.Name)

		slashCommands, err := initSlashCommands()
		if err != nil {
			log.Fatalf("Error initializing slash commands: %s", err.Error())
		}

		if _, err = discordSession.ApplicationCommandBulkOverwrite(discordSession.State.User.ID, i.ID, slashCommands); err != nil {
			log.Fatalf("Error creating slash commands: %s", err.Error())
		}
	})

	discordSession.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.ApplicationCommandData().Name {
		case "team":
			commandHandlers.HandleTeam(session, interaction)
		}
	})

	if err := discordSession.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer discordSession.Close()

	log.Println("Bot is ready")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
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
