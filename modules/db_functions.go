package modules

import "mubasirNew/db"

func verifyAdmin(discordId string) bool {
	var user db.Users
	db.DB.First(&user, "discord_id = ?", discordId)
	if user.IsAdmin {
		return true
	} else {
		return false
	}
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
