package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var componentReply = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"close": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = s.ChannelDelete(i.ChannelID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		delete(TicketChannelMapping, TicketChannelMapping[i.ChannelID])
		delete(TicketChannelMapping, i.ChannelID)

	},
}

func interactionReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := componentReply[i.MessageComponentData().CustomID]; ok {
		h(s, i)
	}
}
