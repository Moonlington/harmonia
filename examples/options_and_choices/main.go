package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Moonlington/harmonia"
	"github.com/bwmarrin/discordgo"
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
	cmd, _ := h.GuildAddSlashCommand("blep", "Send a random adorable animal photo", *GuildID, func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		smol := ""
		if i.GetOption("only_smol") != nil && i.GetOption("only_smol").BoolValue() {
			smol = "baby "
		}
		h.Respond(i, fmt.Sprintf("Sending picture of a %v%v!", smol, i.GetOption("animal").StringValue()))
	})

	opt, _ := cmd.AddOption("animal", "The type of animal", true, discordgo.ApplicationCommandOptionString)
	opt.AddChoice("Dog", "dog")
	opt.AddChoice("Cat", "dog")
	opt.AddChoice("Penguin", "penguin")

	cmd.AddOption("only_smol", "Whether to show only baby animals", false, discordgo.ApplicationCommandOptionBoolean)

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
