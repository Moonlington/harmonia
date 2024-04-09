package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

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
	cmd := harmonia.NewSlashCommand("ping").
		WithDescription("Responds to the user with 'Pong!'... and then again after 5 seconds!").
		WithGuildID(*GuildID).
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			h.Respond(i, "Pong!")

			time.Sleep(time.Second * 5)

			f, _ := h.Followup(i, "Pong again!")

			time.Sleep(time.Second * 5)

			h.EditFollowup(f, "Edited Pong!")
		})

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
