package modules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mubasirNew/db"
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

	switch {
	case m.Content == "!m":
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

	case strings.HasPrefix(m.Content, "!mubasir"):
		msg := strings.Split(m.Content, " ")

		if len(msg) < 3 {
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "<@" + m.Author.ID + "> kayıt için `!mubasir <kullanıcı> <metin>` kullanınız.",
			})
			return
		}

		changeUserId := m.Mentions[0].ID

		if verifyAdmin(m.Author.ID) || m.Author.ID == changeUserId {
			var user db.Users

			customMsg := strings.Join(msg[2:], " ")
			if verifyUser(changeUserId) {
				db.DB.First(&user, "discord_id = ?", changeUserId)
				db.DB.Model(&user).Updates(db.Users{CustomMessage: customMsg})
			} else {
				db.DB.Create(&db.Users{DiscordId: changeUserId, CustomMessage: customMsg, IsAdmin: false})
			}

			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "<@" + changeUserId + "> kullanıcısının karşılama metni `" + customMsg + "` olarak güncellendi!",
			})
		} else {
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "<@" + m.Author.ID + "> kullanıcısı, <@" + changeUserId + "> kullanıcısı için yetkisi olmadığı halde işlem yapmaya çalıştı!",
			})
		}
	}
	return
}
