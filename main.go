package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"dustin-ward/AdventOfCodeBot/bot"
	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

const (
	REQUEST_RATE = time.Minute * 15
)

func main() {
	// Initialize Discord Session
	Session, err := bot.InitSession()
	if err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}

	// Add init handler
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open connection
	if err = Session.Open(); err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}
	defer Session.Close()
	log.Println("Session initialized for", len(bot.C), "servers")

	// Register commands
	r, err := bot.RegisterCommands()
	if err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}
	for _, c := range r {
		log.Printf("Command registered: \"%s\" with id: %v", c.Name, c.ID)
	}

	// Setup Cron
	if err = bot.SetupNotifications(); err != nil {
		log.Println("Error: unable to send notification: %w", err)
	}

	// Continually fetch advent of code data every 15 minutes
	for _, ch := range bot.C {
		go func(channel *data.Channel) {
			for {
				log.Println("Attempting to fetch data for leaderboard " + channel.Leaderboard + "...")
				if err := data.FetchData(channel.Leaderboard, channel.SessionToken, channel.Leaderboard); err != nil {
					log.Println("Error:", fmt.Errorf("fetch: %w", err))
				} else {
					log.Println(channel.Leaderboard, "success!")
				}

				time.Sleep(REQUEST_RATE)
			}
		}(ch)
	}

	// Wait for SIGINT to end program
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	bot.TakeDown()
}
