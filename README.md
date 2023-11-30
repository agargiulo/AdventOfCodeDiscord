# ⭐⭐🎄 Advent Of Code Discord Bot 🎄⭐⭐

A discord bot written in Go to provide leaderboards and nightly notifications to your server!

## Usage 

### General
`/leaderboard` - Show the current state of the leaderboard

### Admin
`/configure <channelId> <roleId> <leaderboardId> <sessionToken>` - Configure the bot to work on your server. Provide the desired channel id, role id to ping, the id of your AoC leaderboard and a valid session token of one member that has access to the private leaderboard

`/start-notifications` - Subscribe to nightly notifications

`/stop-notifications` - Unsubscribe from nightly notifications

`/check-notifications` - Check your servers current notification status

## Deployment
### Using Golang
Ensure you have go installed and run the bot by doing:
```sh
TOKEN="[Your token here]" go run main.go
```

### Using Docker
Ensure you have docker installed and create a volume:
```sh
docker volume create discord-aoc-bot
```
then run the bot by doing:
```sh
docker run -v discord-aoc-bot:/app/ --rm -e TOKEN="[Your token here]" ghcr.io/dustin-ward/advent-of-code-discord:latest
```
