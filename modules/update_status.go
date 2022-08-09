package modules

import (
	"github.com/bwmarrin/discordgo"
)

func UpdateStatus(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "bot")
}
