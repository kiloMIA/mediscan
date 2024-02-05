package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kiloMIA/mediscan/internal"
)

type Bot struct {
    Session       *discordgo.Session
    UserController *internal.UserController
}

func main() {
    token := os.Getenv("DISCORD_BOT_TOKEN")
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("Error creating Discord session,", err)
        return
    }

    db, err := internal.ConnectDB()
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }

    userController := internal.NewUserController(db)
    bot := &internal.Bot{Session: dg, UserController: userController}

    dg.AddHandler(bot.MessageCreate)

    err = dg.Open()
    if err != nil {
        fmt.Println("Error opening connection,", err)
        return
    }

    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

