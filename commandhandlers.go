package harmonia

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is an interface for command handlers.
type CommandHandler interface {
	// Do executes the CommandHandler.
	Do(h *Harmonia, i *Invocation)
	// GetOptions returns the options of the CommandHandler.
	GetOptions() []*discordgo.ApplicationCommandOption
}

// A SingleCommandHandler describes the handler for a command with no subcommands.
type SingleCommandHandler struct {
	// Handler is the actual function of the SlashCommand
	Handler func(h *Harmonia, i *Invocation)
	// Options contain a slice of Options passed through the SlashCommand
	Options []*Option
}

// Do executes the handler.
func (s *SingleCommandHandler) Do(h *Harmonia, i *Invocation) {
	s.Handler(h, i)
}

// GetOptions return the Options given to the SingleCommandHandler.
func (s *SingleCommandHandler) GetOptions() []*discordgo.ApplicationCommandOption {
	o := make([]*discordgo.ApplicationCommandOption, len(s.Options))
	for i, v := range s.Options {
		o[i] = v.ApplicationCommandOption
	}
	return o
}

// A CommandGroupHandler describes the handler for a command with subcommands.
type CommandGroupHandler struct {
	// Subcommands contains a map of SubSlashCommands
	Subcommands map[string]*SlashSubcommand
}

// Do handles the subcommands given to the SubcommandHandler
func (s *CommandGroupHandler) Do(h *Harmonia, i *Invocation) {
	options := i.options
	if sc, ok := s.Subcommands[options[0].Name]; ok {
		i.options = options[0].Options
		sc.Handler.Do(h, i)
	}
}

// GetOptions returns the subcommands parsed as ApplicationCommandOptions
func (s *CommandGroupHandler) GetOptions() []*discordgo.ApplicationCommandOption {
	options := make([]*discordgo.ApplicationCommandOption, len(s.Subcommands))
	i := 0
	for _, sc := range s.Subcommands {
		t := discordgo.ApplicationCommandOptionSubCommand
		if sc.IsGroup {
			t = discordgo.ApplicationCommandOptionSubCommandGroup
		}

		options[i] = &discordgo.ApplicationCommandOption{
			Name:        sc.Name,
			Description: sc.Description,
			Options:     sc.Handler.GetOptions(),
			Type:        t,
		}
		i++
	}
	return options
}
