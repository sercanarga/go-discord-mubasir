package modules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	Buffer = make([][]byte, 0)
)

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

	TextToSpeech(getUser.Username+" geldi! selamlar.", "tmp/output.mp3")
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
