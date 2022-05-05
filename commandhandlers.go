package harmonia

import (
	"github.com/bwmarrin/discordgo"
)

type CommandHandler interface {
	Do(h *Harmonia, i *Invocation)
	GetOptions() []*discordgo.ApplicationCommandOption
}

type SingleCommandHandler struct {
	Handler func(h *Harmonia, i *Invocation)
	Options []*Option
}

func (s *SingleCommandHandler) Do(h *Harmonia, i *Invocation) {
	s.Handler(h, i)
}

func (s *SingleCommandHandler) GetOptions() []*discordgo.ApplicationCommandOption {
	o := make([]*discordgo.ApplicationCommandOption, len(s.Options))
	for i, v := range s.Options {
		o[i] = v.ApplicationCommandOption
	}
	return o
}

type SubcommandHandler struct {
	Subcommands map[string]*SubSlashCommand
}

func (s *SubcommandHandler) Do(h *Harmonia, i *Invocation) {
	options := i.options
	if sc, ok := s.Subcommands[options[0].Name]; ok {
		i.options = options[0].Options
		sc.Handler.Do(h, i)
	}
}

func (s *SubcommandHandler) GetOptions() []*discordgo.ApplicationCommandOption {
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
