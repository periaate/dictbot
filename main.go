package main

import (
	"dictbot/cmd"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DICTBOT_TOKEN")
	guild := os.Getenv("DICTBOT_GUILD")

	Bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	Bot.AddHandler(func(Bot *discordgo.Session, i *discordgo.InteractionCreate) {
		if fn, ok := cmd.Handlers[i.ApplicationCommandData().Name]; ok {
			fn(Bot, i)
		}
	})

	Bot.AddHandler(func(Bot *discordgo.Session, m *discordgo.MessageCreate) {
		if fn, ok := cmd.MessageCreate[m.Content]; ok {
			fn(Bot, m)
		}
	})

	Bot.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", bot.State.User.Username, bot.State.User.Discriminator)
	})

	err = Bot.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// remove all commands
	commands, err := Bot.ApplicationCommands(Bot.State.User.ID, guild)
	if err != nil {
		log.Fatalf("Cannot get commands: %v", err)
	}
	for _, v := range commands {
		err = Bot.ApplicationCommandDelete(Bot.State.User.ID, guild, v.ID)
		if err != nil {
			log.Fatalf("Cannot delete '%v' Command: %v", v.Name, err)
		}
	}

	log.Println("Adding cmd.Commands...")
	for k, v := range cmd.Commands {
		_, err := Bot.ApplicationCommandCreate(Bot.State.User.ID, guild, v)
		if err != nil {
			log.Panicf("Cannot create '%v' Command: %v", v.Name, err)
			// Remove command from map
			delete(cmd.Commands, k)
		}
	}

	defer Bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
