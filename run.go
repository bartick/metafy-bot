package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func runServer(Token string) *discordgo.Session {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}

	dg.AddHandler(MessageCreate)
	dg.AddHandler(interactionReply)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil
	}

	u, err := dg.User("@me")

	if err != nil {
		fmt.Println("Error getting user")
	}

	fmt.Println("Logged in as " + u.Username)

	err = dg.UpdateGameStatus(0, "with your web3")

	if err != nil {
		fmt.Println("Error updating game status")
	}

	return dg
}
