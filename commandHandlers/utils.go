package commandHandlers

import "github.com/bwmarrin/discordgo"

func respondToInteraction(response string, interaction *discordgo.InteractionCreate, session *discordgo.Session) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
