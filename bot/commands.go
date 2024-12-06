package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

func (ab *AocBot) leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"leaderboard\" executed from guildId: %s", i.GuildID)

	// Get calling channel
	channel, err := ab.getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("leaderboard: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}

	// Get leaderboard data for channel
	D, err := data.GetData(channel.Leaderboard)
	if err != nil {
		log.Println("Error:", fmt.Errorf("leaderboard: %v", err))
		respondWithError(s, i, "An internal error occurred...")
		return
	}

	// Sort users by stars and local score
	M := make([]data.User, 0)
	for _, m := range D.Members {
		M = append(M, m)
	}
	sort.Slice(M, func(i, j int) bool {
		if M[i].Stars == M[j].Stars {
			return M[i].LocalScore > M[j].LocalScore
		}
		return M[i].Stars > M[j].Stars
	})

	// Add users to embed
	fields := make([]*discordgo.MessageEmbedField, 0)
	for _, m := range M {
		// Calculate avg. delta time.
		daysFullyComplete := int64(0)
		deltaTimeSum := time.Duration(0)
		for _, d := range m.CompletionDayLevel {
			if d.Silver != nil && d.Gold != nil {
				daysFullyComplete++
				deltaTimeSum += time.Unix(int64(d.Gold.Timestamp), 0).Sub(time.Unix(int64(d.Silver.Timestamp), 0))
			}
		}

		var avgDeltaTime time.Duration = 0
		if daysFullyComplete != 0 {
			avgDeltaTime = time.Duration(deltaTimeSum.Nanoseconds() / daysFullyComplete).Round(time.Millisecond)
		}

		f := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s**", m.Name),
			Value: fmt.Sprintf("⭐ %d 🏆 %d ⏳ %s", m.Stars, m.LocalScore, avgDeltaTime),
		}
		fields = append(fields, f)
	}

	if len(M) > 0 {
		fields[0].Name += " 🥇"
	}
	if len(M) > 1 {
		fields[1].Name += " 🥈"
	}
	if len(M) > 2 {
		fields[2].Name += " 🥉"
	}

	// Create embed object
	embeds := make([]*discordgo.MessageEmbed, 1)
	embeds[0] = &discordgo.MessageEmbed{
		URL:   data.AocLeaderboardUrl(channel.Leaderboard),
		Type:  discordgo.EmbedTypeRich,
		Title: data.LearderBoardTitle,
		Color: 0x127C06,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "⏳: Average time to solve part 2",
		},
		Description: fmt.Sprintf("Leaderboard as of <t:%d:F>", time.Now().Unix()),
		Fields:      fields,
	}

	// Send embed to channel
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
	if err != nil {
		log.Println("Warn:", fmt.Errorf("leaderboard: %w", err))
	}
}

func (ab *AocBot) configure(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"configure\" executed from guildId: %s", i.GuildID)

	// Grab command options from user
	options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		options[opt.Name] = opt
	}

	// Create new channel object
	ch := data.Channel{
		GuildId:         i.GuildID,
		ChannelId:       options["channel"].ChannelValue(nil).ID,
		RoleId:          options["role"].RoleValue(nil, i.GuildID).ID,
		Leaderboard:     options["leaderboard"].StringValue(),
		SessionToken:    options["session-token"].StringValue(),
		NotificationsOn: false,
	}

	// Add to local memory
	ab.c[i.GuildID] = &ch

	// Write to file
	b, err := json.Marshal(ab.c)
	if err != nil {
		log.Println("Error:", fmt.Errorf("configure: %v", err))
		respondWithError(s, i, "Error: Invalid arguments were supplied...")
		return
	}

	if err := os.WriteFile(dataDir+"channels.json", b, 0777); err != nil {
		log.Println("Error:", fmt.Errorf("configure: %v", err))
		respondWithError(s, i, "Error: Internal server error...")
		return
	}

	log.Println("Attempting to fetch data for leaderboard " + ch.Leaderboard + "...")
	if err := data.FetchData(ch.Leaderboard, ch.SessionToken, ch.Leaderboard); err != nil {
		log.Println("Error:", fmt.Errorf("fetch: %w", err))
	} else {
		log.Println(ch.Leaderboard, "success!")
	}

	respond(s, i, "Server successfully configured!", true)
}

func (ab *AocBot) startCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"start-notifications\" executed from guildId: %s", i.GuildID)

	ch, err := ab.getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("start-notifications: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}
	ch.NotificationsOn = true

	respond(s, i, "Notification process started! ⏰", false)
}

func (ab *AocBot) stopCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"stop-notifications\" executed from guildId: %s", i.GuildID)

	ch, err := ab.getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("start-notifications: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}
	ch.NotificationsOn = false

	respond(s, i, "Notification process stopped! ⏸", false)
}

func (ab *AocBot) checkCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"check-notifications\" executed from guildId: %s", i.GuildID)
	ch, err := ab.getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("check-notifications: %w", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}

	next, err := ab.nextNotification()
	if err != nil {
		log.Println("Error:", fmt.Errorf("check-notifications: %w", err))
		respondWithError(s, i, "Internal Error 💀 Did you configure the notifications correctly?")
		return
	}
	day := next.Day()

	var message string
	if ch.NotificationsOn {
		message = fmt.Sprintf("Notifications for server id: %s are enabled in channel: %s!\n\n⏰ Next notification: <t:%d:R> (Day %d)", ch.GuildId, ch.ChannelId, next.Unix(), day)
	} else {
		message = "Notifications are not enabled currently..."
	}

	respond(s, i, message, false)
}
