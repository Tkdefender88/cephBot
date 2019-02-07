package bot

type guilds struct {
	Count  int `json:"server_count"`
	Server map[string]*guild
}

type guild struct {
	GuildID       string           `json:"guildID"`
	CommandPrefix string           `json:"prefix"`
	EmbedColor    int              `json:"embed_color"`
	Kicked        bool             `json:"kicked"`
	LimitStars    bool             `json:"limStars"`
	Users         map[string]*user `json:"users"`
}

type user struct {
	UserID string `json:"userID"`
	Stars  int    `json:"stars"`
}
