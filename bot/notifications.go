package bot

import (
	"fmt"
	"log"
	"time"

	"dustin-ward/AdventOfCodeBot/data"
)

func (ab *AocBot) problemNotification() {
	day := time.Now().In(time.UTC).Day()
	// Okay please stop notifying outside of Advent of Code.
	if now := time.Now().In(time.UTC); !(now.Month() == time.December && now.Day() <= 25) {
		return
	}
	// For each registered channel
	for _, ch := range ab.c {
		if ch.NotificationsOn {
			log.Println("Info: sending day", day, "notification in channel", ch.ChannelId)

			// Create message object
			messageString := fmt.Sprintf(
				"🎄 <@&%s> 🎄\nThe problem for Day %d will be released soon! (<t:%d:R>)\nYou can see the problem statement here when its up: %s",
				ch.RoleId,
				day,
				time.Now().Add(time.Minute*30).Unix(),
				data.ProblemUrl(day),
			)

			// Send message to channel
			_, err := ab.s.ChannelMessageSend(ch.ChannelId, messageString)
			if err != nil {
				log.Println("Error:", fmt.Errorf("unable to send notification: %w", err))
			}
		} else {
			log.Println("Info: notifications disabled for", ch.GuildId)
		}
	}
}

func (ab *AocBot) nextNotification() (time.Time, error) {
	entries := ab.cron.Entries()
	if len(entries) != 1 {
		return time.Now(), fmt.Errorf("invalid number of cron entries")
	}
	return (*entries[0]).Next, nil
}
