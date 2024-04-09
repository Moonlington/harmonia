package main

import (
	"flag"
	"fmt"
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
	h.AddCommand(harmonia.NewUserCommand("Ping this guy.").
		WithGuildID(*GuildID).
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			target, _ := i.TargetAuthor(h)
			h.Respond(i, fmt.Sprintf("Ping %s!", target.Mention()))
		}))

	h.AddCommand(harmonia.NewMessageCommand("Who sent this message?").
		WithGuildID(*GuildID).
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			target, _ := i.TargetMessage(h)
			h.Respond(i, fmt.Sprintf("%s wrote this!", target.Author.Mention()))
		}))

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
