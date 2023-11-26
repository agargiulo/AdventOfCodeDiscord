package bot

import (
	"fmt"
	"log"
	"time"
)

const (
	ProblemUrl = "https://adventofcode.com/2023/day/"
)

func problemNotification() {
	day := time.Now().AddDate(0, 0, 1).Day()

	// For each registered channel
	for _, ch := range C {
		if ch.NotificationsOn {
			log.Println("Info: sending day", day, "notification in channel", ch.ChannelId)

			// Create message object
			messageString := fmt.Sprintf(
				"🎄 <@&%s> 🎄\nThe problem for Day %d will be released soon! (<t:%d:R>)\nYou can see the problem statement here when its up: https://adventofcode.com/2023/day/%d",
				ch.RoleId,
				day,
				time.Now().Unix()+(int64(30)*60),
				day,
			)

			// Send message to channel
			_, err := s.ChannelMessageSend(ch.ChannelId, messageString)
			if err != nil {
				log.Println("Error:", fmt.Errorf("unable to send notification: %w", err))
			}
		} else {
			log.Println("Info: notifications disabled for", ch.GuildId)
		}
	}
}

func NextNotification() (time.Time, error) {
	entries := crn.Entries()
	if len(entries) != 1 {
		return time.Now(), fmt.Errorf("invalid number of cron entries")
	}
	return (*entries[0]).Next, nil
}
