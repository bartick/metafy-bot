package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {

	client := connectDb()

	dg := runServer(os.Getenv("TOKEN"))

	if dg == nil {
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C (^C) to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	client.Disconnect(ctx)
	dg.Close()
}
