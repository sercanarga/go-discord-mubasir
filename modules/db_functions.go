package modules

import (
	"github.com/bwmarrin/discordgo"
	"mubasirv2/db"
)

func verifyAdmin(s *discordgo.Session, discordId string, guildId string) bool {
	guildUser, _ := s.GuildMember(guildId, discordId)

	for _, role := range guildUser.Roles {
		if role == "964335318482452530" {
			return true
		}
	}

	return false
}

func verifyUser(discordId string) bool {
	var user db.Users
	db.DB.First(&user, "discord_id = ?", discordId)
	if user.DiscordId == discordId {
		return true
	} else {
		return false
	}
}
