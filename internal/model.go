package internal

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Appointment struct {
    ID        int64  `json:"id"`
    DiscordID string `json:"discord_id"`
    Timestamp time.Time `json:"timestamp"` 
    Details   string `json:"details"`
}

type Document struct {
    ID       int64  `json:"id"`
    DiscordID string `json:"discord_id"`
    FilePath string `json:"file_path"`
}

type Bot struct {
    Session       *discordgo.Session
    UserController *UserController
}