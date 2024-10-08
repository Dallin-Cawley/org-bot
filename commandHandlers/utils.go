package commandHandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"orgBot/database"
	"orgBot/database/model"
	"strings"
)

func respondToInteraction(response string, interaction *discordgo.InteractionCreate, session *discordgo.Session) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func caseInsensitiveEquality(lhs string, rhs string) bool {
	return strings.ToLower(lhs) == strings.ToLower(rhs)
}

func getBotRoleID(guildID string, transaction pgx.Tx) (string, error) {
	guilds, err := database.ReadModel[model.Guild](model.MakeGuild(guildID, ""), transaction)
	if err != nil {
		return "", err
	}

	return guilds[0].BotRoleID, nil
}

func numMembersInVoiceChannel(channelID string, guildID string, session *discordgo.Session) (int, error) {
	guild, err := session.State.Guild(guildID)
	if err != nil {
		return -1, err
	}

	numMembers := 0
	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == channelID {
			numMembers++
		}
	}

	return numMembers, nil
}
