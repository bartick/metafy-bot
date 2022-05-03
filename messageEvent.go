package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	TicketChannelMapping = map[string]string{}
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel")
		return
	}
	if channel.Type == discordgo.ChannelTypeDM {
		if TicketChannelMapping[m.Author.ID] == "" {
			channel, err := s.GuildChannelCreate(os.Getenv("GUILD_ID"), m.Author.Username, discordgo.ChannelTypeGuildText)
			if err != nil {
				fmt.Println("Error creating channel")
				return
			}

			_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
				ParentID: os.Getenv("CATEGORY_ID"),
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   os.Getenv("GUILD_ID"),
						Deny: discordgo.PermissionViewChannel,
					},
					{
						ID:    os.Getenv("HELPER_ROLE"),
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionViewChannel,
					},
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			TicketChannelMapping[m.Author.ID] = channel.ID
			TicketChannelMapping[channel.ID] = m.Author.ID

			_, err = s.ChannelMessageSend(m.ChannelID, "Hello, **"+m.Author.Username+"!**\nThank you for reaching to us. Please wait and a staff member will be with you shortly.\nIf you have any questions, please feel free to ask them to me.")
			if err != nil {
				fmt.Println("Error sending message")
				return
			}

			_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
				Content: "@here",
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "New Ticket",
						Description: "Seems like **" + m.Author.Username + "** need some help. Please use !msg to reply to them.",
						Color:       0x00ff00,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: s.State.User.AvatarURL(""),
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name:    m.Author.Username,
							IconURL: m.Author.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Close",
								Style:    discordgo.PrimaryButton,
								CustomID: "close",
								Emoji: discordgo.ComponentEmoji{
									Name: "🔒",
								},
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}

		}
		_, err := s.ChannelMessageSend(TicketChannelMapping[m.Author.ID], "**"+m.Author.Username+"#"+m.Author.Discriminator+"**: "+m.Content)
		if err != nil {
			fmt.Println("Error sending message")
			return
		}

	} else {
		re := regexp.MustCompile("\\s+")
		contentArray := re.Split(m.Content, -1)
		if contentArray[0] == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")
		} else if contentArray[0] == "!msg" {
			if len(contentArray) < 2 {
				s.ChannelMessageSend(m.ChannelID, "Please provide a message to send to the user.")
				return
			} else {

				channel, err = s.UserChannelCreate(TicketChannelMapping[m.ChannelID])
				if err != nil {
					fmt.Println("Error creating channel")
					return
				}
				message := strings.Join(contentArray[1:], " ")
				_, err := s.ChannelMessageSend(channel.ID, "**"+m.Author.Username+"**: "+message)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				s.ChannelMessageDelete(m.ChannelID, m.ID)

				_, err = s.ChannelMessageSend(m.ChannelID, "**"+m.Author.Username+"#"+m.Author.Discriminator+"**: "+message)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		} else {
			return
		}
	}

}
