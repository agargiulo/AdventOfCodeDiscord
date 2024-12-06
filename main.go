package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"dustin-ward/AdventOfCodeBot/bot"
	"dustin-ward/AdventOfCodeBot/data"

	"github.com/bwmarrin/discordgo"
)

const (
	RequestRate = time.Minute * 15
)

func main() {
	// Setup the default values for the environment variables
	if os.Getenv("DATA_DIR") == "" {
		if err := os.Setenv("DATA_DIR", "./"); err != nil {
			panic(err)
		}
	}

	// Avoid most user indused errors due to an incorrect path.
	if !strings.HasSuffix(os.Getenv("DATA_DIR"), "/") {
		if err := os.Setenv("DATA_DIR", os.Getenv("DATA_DIR")+"/"); err != nil {
			panic(err)
		}
	}

	// Initialize Discord Session
	aocBot, err := bot.InitBot()
	if err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}

	// Add init handler
	aocBot.Session().AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open connection
	if err = aocBot.Session().Open(); err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}
	defer func(session *discordgo.Session) {
		err := session.Close()
		if err != nil {
			panic(err)
		}
	}(aocBot.Session())
	log.Println("Session initialized for", len(aocBot.Chans()), "servers")

	// Register commands
	r, err := aocBot.RegisterCommands()
	if err != nil {
		log.Fatal("Fatal:", fmt.Errorf("main: %w", err))
	}
	for _, c := range r {
		log.Printf("Command registered: \"%s\" with id: %v", c.Name, c.ID)
	}

	// Setup Cron
	if err = aocBot.SetupNotifications(); err != nil {
		log.Println("Error: unable to send notification: %w", err)
	}

	// Continually fetch advent of code data every 15 minutes
	for _, ch := range aocBot.Chans() {
		go func(channel *data.Channel) {
			for {
				log.Println("Attempting to fetch data for leaderboard " + channel.Leaderboard + "...")
				if err := data.FetchData(channel.Leaderboard, channel.SessionToken, channel.Leaderboard); err != nil {
					log.Println("Error:", fmt.Errorf("fetch: %w", err))
				} else {
					log.Println(channel.Leaderboard, "success!")
				}

				time.Sleep(RequestRate)
			}
		}(ch)
	}

	// Wait for SIGINT to end program
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if err := aocBot.TakeDown(); err != nil {
		panic(err)
	}
}
