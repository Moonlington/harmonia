package harmonia

import (
	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(h *Harmonia, i *Invocation)

type CommandHandler interface {
	// GetName returns the name of the Command.
	GetName() string

	// Do executes the Command.
	Do(h *Harmonia, i *Invocation)

	getRegistration() *discordgo.ApplicationCommand
	setRegistration(*discordgo.ApplicationCommand)
}
