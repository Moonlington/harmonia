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
	opt := harmonia.NewOption("animal", discordgo.ApplicationCommandOptionString).
		WithDescription("The type of animal").
		IsRequired().
		AddChoice("Dog", "dog").
		AddChoice("Cat", "dog").
		AddChoice("Penguin", "penguin")

	opt2 := harmonia.NewOption("only_smol", discordgo.ApplicationCommandOptionBoolean).
		WithDescription("Wether to show only baby animals")

	cmd := harmonia.NewSlashCommand("blep").
		WithDescription("Send an adorable animal photo").
		WithGuildID(*GuildID).
		WithOptions(opt, opt2).
		WithCommand(func(h *harmonia.Harmonia, i *harmonia.Invocation) {
			smol := ""
			if i.GetOption("only_smol") != nil && i.GetOption("only_smol").BoolValue() {
				smol = "baby "
			}
			h.Respond(i, fmt.Sprintf("Sending picture of a %v%v!", smol, i.GetOption("animal").StringValue()))
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
