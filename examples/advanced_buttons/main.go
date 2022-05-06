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
	number := 0

	h.AddSlashCommandInGuild("number", "Increase or decrease the internal number!", *GuildID, func(h *harmonia.Harmonia, i *harmonia.Invocation) {
		h.DeferRespond(i)
		f, err := h.FollowupWithComponents(i, fmt.Sprintf("The current number is %v!", number), [][]discordgo.MessageComponent{
			{
				discordgo.Button{
					Label:    "Increase by 1",
					Style:    discordgo.SuccessButton,
					CustomID: "n_increase",
				},
				discordgo.Button{
					Label:    "Decrease by 1",
					Style:    discordgo.DangerButton,
					CustomID: "n_decrease",
				},
			}, {
				discordgo.Button{
					Label:    "Reset to 0",
					Style:    discordgo.PrimaryButton,
					CustomID: "n_reset",
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		h.AddComponentHandlerToFollowup(f, "n_increase", func(h *harmonia.Harmonia, ci *harmonia.Invocation) {
			if i.Author.ID == ci.Author.ID {
				number++
				h.EphemeralRespond(ci, fmt.Sprintf("The number has been increased to %v", number))
			} else {
				h.EphemeralRespond(ci, "Only the original caller of the function can use it!")
			}
		})
		h.AddComponentHandlerToFollowup(f, "n_decrease", func(h *harmonia.Harmonia, ci *harmonia.Invocation) {
			if i.Author.ID == ci.Author.ID {
				number--
				h.EphemeralRespond(ci, fmt.Sprintf("The number has been decreased to %v", number))
			} else {
				h.EphemeralRespond(ci, "Only the original caller of the function can use it!")
			}
		})
		h.AddComponentHandlerToFollowup(f, "n_reset", func(h *harmonia.Harmonia, ci *harmonia.Invocation) {
			if i.Author.ID == ci.Author.ID {
				number = 0
				h.EphemeralRespond(ci, fmt.Sprintf("The number has been reset to %v", number))
			} else {
				h.EphemeralRespond(ci, "Only the original caller of the function can use it!")
			}
		})
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
