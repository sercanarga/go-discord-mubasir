package modules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var (
	Vc *discordgo.VoiceConnection
)

func joinServer(s *discordgo.Session, guildID, channelID string) error {
	var err error

	Vc, err = s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	return nil
}

func CallBot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!m") {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = joinServer(s, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error coming channel:", err)
				}
				return
			}
		}
	}

	return
}
