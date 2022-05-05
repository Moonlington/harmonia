package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Moonlington/harmonia"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var h *harmonia.Harmonia

func init() { flag.Parse() }

func init() {
	var err error
	h, err = harmonia.New(*BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	c, _ := h.AddSlashCommandWithSubcommandsInGuild("main", "Subcommands example", *GuildID)
	c.AddSubcommand("first", "First subcommand!", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "This is the first subcommand!")
	})
	c.AddSubcommand("second", "Second subcommand!", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "This is the second subcommand!")
	})
	group, _ := c.AddSubcommandGroup("third", "Third subcommand, but a group!")
	group.AddSubcommand("fourth", "Fourth subcommand, but nested in the third one! I guess it's the actual third subcommand?", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "This is the REAL third subcommand!")
	})

	err := h.Run()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer h.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		err := h.RemoveAllCommands()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Gracefully shutting down.")
}
