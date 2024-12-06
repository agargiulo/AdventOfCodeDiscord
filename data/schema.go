package data

type Channel struct {
	GuildId         string `json:"guildId"`
	ChannelId       string `json:"channelId"`
	RoleId          string `json:"roleId"`
	Leaderboard     string `json:"leaderboard"`
	SessionToken    string `json:"token"`
	NotificationsOn bool   `json:"notificationsOn"`
}

type Channels map[string]*Channel

type Data struct {
	Event   string          `json:"event"`
	OwnerID int32           `json:"owner_id"`
	Members map[string]User `json:"members"`
}

type User struct {
	Id                 uint32         `json:"id"`
	Name               string         `json:"name"`
	GlobalScore        uint32         `json:"global_score"`
	LocalScore         uint32         `json:"local_score"`
	Stars              uint32         `json:"stars"`
	LastStar           uint32         `json:"last_star_ts"`
	CompletionDayLevel map[string]Day `json:"completion_day_level"`
}

type Day struct {
	Silver *Star `json:"1"`
	Gold   *Star `json:"2"`
}

type Star struct {
	Index     uint32 `json:"star_index"`
	Timestamp uint32 `json:"get_star_ts"`
}
