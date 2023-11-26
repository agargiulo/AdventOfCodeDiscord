package bot

import (
	"dustin-ward/AdventOfCodeBot/data"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

func helloworld(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"helloworld\" executed from guildId: %s", i.GuildID)

	respond(s, i, "Hello from the AoC bot 🙂")
}

func leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"leaderboard\" executed from guildId: %s", i.GuildID)

	// Get calling channel
	channel, err := getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("leaderboard: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}

	// Get leaderboard data for channel
	D, err := data.GetData(channel.Leaderboard)
	if err != nil {
		log.Println("Error:", fmt.Errorf("leaderboard: %v", err))
		respondWithError(s, i, "An internal error occured...")
		return
	}

	// Sort users by stars and localscore
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
		daysFullyComplete := uint32(0)
		deltaTimeSum := uint32(0)
		for _, d := range m.CompletionDayLevel {
			if d.Silver != nil && d.Gold != nil {
				daysFullyComplete++
				deltaTimeSum += d.Gold.Timestamp - d.Silver.Timestamp
			}
		}

		var avgDeltaTime uint32 = 0
		var avgDtimeS uint32 = 0
		var avgDtimeM uint32 = 0
		if daysFullyComplete != 0 {
			avgDeltaTime = deltaTimeSum / daysFullyComplete
			avgDtimeS = avgDeltaTime % 60
			avgDtimeM = avgDeltaTime / 60
		}

		f := &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s**", m.Name),
			Value: fmt.Sprintf("⭐ %d 🏆 %d ⏳ %d:%02d", m.Stars, m.LocalScore, avgDtimeM, avgDtimeS),
		}
		fields = append(fields, f)
	}
	fields[0].Name += " 🥇"
	fields[1].Name += " 🥈"
	fields[2].Name += " 🥉"

	// Create embed object
	embeds := make([]*discordgo.MessageEmbed, 1)
	embeds[0] = &discordgo.MessageEmbed{
		URL:   "https://adventofcode.com/2022/private/view/" + channel.Leaderboard,
		Type:  discordgo.EmbedTypeRich,
		Title: "🎄 2022 Leaderboard 🎄",
		Color: 0x127C06,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Leaderboard as of " + time.Now().Format("2006/01/02 3:4:5pm"),
		},
		Fields: fields,
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

func configure(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	C[i.GuildID] = &ch

	// Write to file
	b, err := json.Marshal(C)
	if err != nil {
		log.Println("Error:", fmt.Errorf("configure: %v", err))
		respondWithError(s, i, "Error: Invalid arguments were supplied...")
		return
	}

	if err := os.WriteFile("./channels.json", b, 0777); err != nil {
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

	respond(s, i, "Server successfully configured!")
}

func startCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"start-notifications\" executed from guildId: %s", i.GuildID)

	ch, err := getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("start-notifications: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}
	ch.NotificationsOn = true

	respond(s, i, "Notification process started! ⏰")
}

func stopCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"stop-notifications\" executed from guildId: %s", i.GuildID)

	ch, err := getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("start-notifications: %v", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}
	ch.NotificationsOn = false

	respond(s, i, "Notification process stopped! ⏸")
}

func checkCountdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Info: command \"check-notifications\" executed from guildId: %s", i.GuildID)
	ch, err := getChannel(i.GuildID)
	if err != nil {
		log.Println("Error:", fmt.Errorf("check-notifications: %w", err))
		respondWithError(s, i, "Your server has not been correctly configured! 🛠️ Use /configure-server")
		return
	}

	next, err := NextNotification()
	if err != nil {
		log.Println("Error:", fmt.Errorf("check-notifications: %w", err))
		respondWithError(s, i, "Internal Error 💀 Please contact @shrublord")
		return
	}
	day := time.Now().AddDate(0, 0, 1).Day()

	var message string
	if ch.NotificationsOn {
		message = fmt.Sprintf("Notifications for server id: %s are enabled in channel: %s!\n\n⏰ Next notification: <t:%d:R> (Day %d)", ch.GuildId, ch.ChannelId, next.Unix(), day)
	} else {
		message = "Notifications are not enabled currently..."
	}

	respond(s, i, message)
}
