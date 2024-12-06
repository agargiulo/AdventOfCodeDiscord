package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	adminPerm int64 = 0
	dataDir         = os.Getenv("DATA_DIR")
)

type commandFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type AocBot struct {
	s        *discordgo.Session
	c        data.Channels
	cron     *cron.Cron
	commands map[string]commandFunc
}

func InitBot() (*AocBot, error) {
	ab := &AocBot{}
	ab.commands = map[string]commandFunc{
		"leaderboard":         ab.leaderboard,
		"configure-server":    ab.configure,
		"start-notifications": ab.startCountdown,
		"stop-notifications":  ab.stopCountdown,
		"check-notifications": ab.checkCountdown,
	}

	if err := ab.initSession(); err != nil {
		return nil, err
	}
	return ab, nil
}

func (ab *AocBot) Session() *discordgo.Session {
	return ab.s
}

func (ab *AocBot) Chans() data.Channels {
	return ab.c
}

func (ab *AocBot) initSession() error {
	// Get token
	token := os.Getenv("AOC_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("no discord token found. Please set $AOC_BOT_TOKEN")
	}

	// Init discordgo session
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("invalid bot configuration: %w", err)
	}

	// Attach handlers to functions
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := ab.commands[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		} else {
			log.Printf("Warning: command handler not found: \"%s\"\n", i.ApplicationCommandData().Name)
		}
	})

	// Check for config file
	if _, err := os.Stat(dataDir + "channels.json"); errors.Is(err, os.ErrNotExist) {
		log.Println("Info: no channel config file found")

		ab.c = make(map[string]*data.Channel, 3)
	} else {
		// Read channel configs from file (Not an ideal storage method...)
		b, err := os.ReadFile(dataDir + "channels.json")
		if err != nil {
			return err
		}

		// Populate channel info in local memory
		if err = json.Unmarshal(b, &ab.c); err != nil {
			return err
		}
	}

	ab.s = s
	return nil
}

func (ab *AocBot) TakeDown() error {
	log.Println("Shutting Down...")
	ab.cron.Stop()

	// Save channel configurations
	b, err := json.Marshal(ab.c)
	if err != nil {
		return err
	}
	if err = os.WriteFile(dataDir+"channels.json", b, 0777); err != nil {
		return err
	}

	return nil
}

func (ab *AocBot) RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := ab.s.ApplicationCommandCreate(ab.s.State.User.ID, "", v)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%s' command: %w", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return registeredCommands, nil
}

func (ab *AocBot) SetupNotifications() error {
	// Cronjob for 11:30pm EST (04:30 UTC)
	ab.cron = cron.NewWithLocation(time.UTC)
	if err := ab.cron.AddFunc("0 30 4 * * *", ab.problemNotification); err != nil {
		return err
	}
	ab.cron.Start()
	return nil
}

// Command definitions
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "leaderboard",
		Description: "Current Leaderboard",
	},
	{
		Name:                     "configure-server",
		Description:              "Configure the AdventOfCode bot to work with your leaderboard and server",
		DefaultMemberPermissions: &adminPerm,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Channel to post in",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Advent of Code role to be pinged",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "leaderboard",
				Description: "Id for your private Advent of Code leaderboard",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "session-token",
				Description: "Valid session token of one member of your private leaderboard",
				Required:    true,
			},
		},
	},
	{
		Name:                     "start-notifications",
		Description:              "Start the notification process",
		DefaultMemberPermissions: &adminPerm,
	},
	{
		Name:                     "stop-notifications",
		Description:              "Stop the notification process",
		DefaultMemberPermissions: &adminPerm,
	},
	{
		Name:                     "check-notifications",
		Description:              "Check to see if notifications are currently enabled",
		DefaultMemberPermissions: &adminPerm,
	},
}
