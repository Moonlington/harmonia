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
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdown or not")
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
	cmd := harmonia.NewGroupSlashCommand("main").
		WithDescription("Subcommands example").
		WithGuildID(*GuildID)

	cmd1 := harmonia.NewSlashCommand("first").
		WithDescription("First subcommand!").
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			h.Respond(i, "This is the first subcommand!")
		})

	cmd2 := harmonia.NewSlashCommand("second").
		WithDescription("Second subcommand!").
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			h.Respond(i, "This is the second subcommand!")
		})

	cmd3 := harmonia.NewGroupSlashCommand("third").
		WithDescription("Third subcommand, but a group!")

	cmd4 := harmonia.NewSlashCommand("fourth").
		WithDescription("Fourth subcommand, but nested in the third one! I guess it's the actual third subcommand?").
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			h.Respond(i, "This is the REAL third subcommand!")
		})

	cmd3.WithSubCommands(cmd4)

	cmd.WithSubCommands(cmd1, cmd2, cmd3)

	h.AddCommand(cmd)

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
