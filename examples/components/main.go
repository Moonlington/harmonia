package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

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
	h.GuildAddSlashCommand("buttons", "Gives you buttons to press!", *GuildID, func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.RespondWithComponents(i, "Look at all these buttons!", [][]discordgo.MessageComponent{
			{
				discordgo.Button{Label: "Say yes!", CustomID: "button1"},
				discordgo.Button{Label: "Say no!", CustomID: "button2"},
			}, {
				discordgo.Button{Label: "Say... maybe?", CustomID: "button3"},
			},
		})
	})

	h.AddComponentHandler("button1", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "Yes!")
	})
	h.AddComponentHandler("button2", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "No!")
	})
	h.AddComponentHandler("button3", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "What?")
	})

	h.GuildAddSlashCommand("class", "Gives you a selection of classes!", *GuildID, func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		min := 1
		h.RespondWithComponents(i, "Choose a couple of classes!", [][]discordgo.MessageComponent{
			{
				discordgo.SelectMenu{
					CustomID:    "class_select_1",
					Placeholder: "Choose a class",
					MinValues:   &min,
					MaxValues:   3,
					Options: []discordgo.SelectMenuOption{
						{
							Label:       "Rogue",
							Value:       "rogue",
							Description: "Sneak n stab",
						},
						{
							Label:       "Mage",
							Value:       "mage",
							Description: "Turn 'em into a sheep",
						},
						{
							Label:       "Priest",
							Value:       "priest",
							Description: "You get heals when I'm done doing damage",
						},
					},
				},
			},
		})
	})

	h.AddComponentHandler("class_select_1", func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.Respond(i, "You have responded with "+strings.Join(i.Values, ", "))
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
