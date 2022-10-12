package modules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mubasirv2/db"
	"os"
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
	case m.Content == "!mubasir":
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

	case strings.HasPrefix(m.Content, "!m"):
		msg := strings.Split(m.Content, " ")

		if m.Attachments == nil && (len(msg) < 3 || len(m.Mentions) != 1) {
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "<@" + m.Author.ID + "> kayıt için `!m <kullanıcı> <metin>` kullanınız.",
			})
			return
		}

		changeUserId := ""
		if m.Mentions[0] != nil {
			changeUserId = m.Mentions[0].ID
		} else {
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "<@" + m.Author.ID + "> kayıt için `!m <kullanıcı> <metin>` kullanınız.",
			})
			return
		}

		if verifyAdmin(s, m.Author.ID, m.GuildID) || m.Author.ID == changeUserId {
			var user db.Users

			if len(m.Attachments) == 1 {
				url := strings.Split(m.Attachments[0].URL, ".")
				extension := url[len(url)-1]

				if extension == "mp3" && m.Attachments[0].Size < 1000000 {
					err := DownloadFile("tmp/"+changeUserId+".mp3", m.Attachments[0].URL)
					if err != nil {
						_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
							Content: "Dosya discord CDN'den indirilemedi!",
						})
						return
					}

					err = ConvertDCA("tmp/"+changeUserId+".mp3", "tmp/"+changeUserId+".dca")
					_ = os.Remove("tmp/" + changeUserId + ".mp3")
					if err != nil {
						_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
							Content: "Ses dosyası dönüştürülemedi!",
						})
						return
					}
				} else {
					_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
						Content: "En fazla 1MB boyutunda `.mp3` uzantılı dosyalar yüklenebilir.",
					})
					return
				}
				_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content: "<@" + changeUserId + "> kullanıcısının karşılama metni `SES` olarak güncellendi!",
				})
				return
			}
			if _, err := os.Stat("tmp/" + changeUserId + ".dca"); err == nil {
				_ = os.Remove("tmp/" + changeUserId + ".dca")
			}

			customMsg := strings.Join(msg[2:], " ")
			if len(customMsg) > 50 {
				customMsg = ""
			}

			if verifyUser(changeUserId) {
				db.DB.First(&user, "discord_id = ?", changeUserId)
				db.DB.Model(&user).Updates(db.Users{CustomMessage: customMsg})
			} else {
				db.DB.Create(&db.Users{DiscordId: changeUserId, CustomMessage: customMsg})
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
