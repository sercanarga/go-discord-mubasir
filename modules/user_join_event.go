package modules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mubasirv2/db"
	"regexp"
)

var (
	Buffer               = make([][]byte, 0)
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

// BotChangeChannelEvent @todo: geliştirilecek
func BotChangeChannelEvent(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.BeforeUpdate != nil {
		if vsu.BeforeUpdate.ChannelID != "" {
			if vsu.BeforeUpdate.UserID == s.State.User.ID {
				joinServer(s, vsu.GuildID, vsu.ChannelID)
				fmt.Println("Bot changed channel")
			}
		}
	}
}

func UserJoinEvent(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	Buffer = make([][]byte, 0)
	getUser, _ := s.User(vsu.UserID)

	if getUser.ID == s.State.User.ID {
		return
	}

	if vsu.BeforeUpdate != nil {
		if vsu.BeforeUpdate.ChannelID != "" {
			return
		}
	}

	var user db.Users
	db.DB.Where("discord_id = ?", vsu.UserID).First(&user)

	var joinUsername = clearString(getUser.Username)
	if len(joinUsername) > 15 {
		joinUsername = clearString(joinUsername[:15])
	}

	if user.DiscordId != "" {
		TextToSpeech(joinUsername+" geldi! "+user.CustomMessage, "tmp/output.mp3")
	} else {
		TextToSpeech(joinUsername+" geldi! Kim bu tanımıyorum?", "tmp/output.mp3")
	}

	ConvertDCA("tmp/output.mp3", "tmp/output.dca")
	err := loadSound()

	if err != nil {
		fmt.Println("Error loading sound: ", err)
	}

	err = PlaySound(Vc)
	if err != nil {
		fmt.Println(err)
	}
}
