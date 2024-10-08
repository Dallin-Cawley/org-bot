package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"orgBot/commandHandlers"
	"orgBot/database"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := database.Connect(context.Background(), os.Getenv("DB_DSN")); err != nil {
		log.Fatal(err)
		return
	}

	discordSession, err := discordgo.New("Bot " + os.Getenv("DISCORD_API_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	discordSession.Identify.Intents = discordgo.IntentsAll

	discordSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	discordSession.AddHandler(commandHandlers.GuildCreate)
	discordSession.AddHandler(commandHandlers.VoiceChannelStatusUpdate)

	discordSession.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.ApplicationCommandData().Name {
		case "team":
			commandHandlers.HandleTeam(session, interaction)
		case "join_to_create":
			commandHandlers.JoinToCreateVoiceChat(session, interaction)
		case "copy_role":
			commandHandlers.AddRoleSync(session, interaction)
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
