package bot

import (
	"dustin-ward/AdventOfCodeBot/data"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func getChannel(guildId string) (*data.Channel, error) {
	ch, ok := C[guildId]
	if !ok {
		return nil, fmt.Errorf("channel not found")
	}
	return ch, nil
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
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
