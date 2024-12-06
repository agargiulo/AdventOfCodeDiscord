package bot

import (
	"fmt"
	"log"

	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

func (ab *AocBot) getChannel(guildId string) (*data.Channel, error) {
	ch, ok := ab.c[guildId]
	if !ok {
		return nil, fmt.Errorf("channel not found")
	}
	return ch, nil
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, message string, ephemeral bool) {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags |= discordgo.MessageFlagsEphemeral
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   flags,
		},
	})
	if err != nil {
		log.Println("Warn:", fmt.Errorf("stop-notifications: %w", err))
	}
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		log.Fatal("Fatal:", fmt.Errorf("respondWithError: %v", err))
	}
}
