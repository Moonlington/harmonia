package harmonia

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Name         string
	Description  string
	GuildID      string
	Handler      CommandHandler
	registration *discordgo.ApplicationCommand
}

func (s *SlashCommand) AddSubcommand(name, description string, handler func(h *Harmonia, i *Invocation)) (sub *SubSlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Subcommand name")
	}

	ch, ok := s.Handler.(*SubcommandHandler)
	if !ok {
		return nil, fmt.Errorf("Slash Command '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("Subcommand '%v' already exists", name)
	}

	sub = &SubSlashCommand{
		Name:        name,
		Description: description,
		Handler:     &SingleCommandHandler{Handler: handler},
	}

	ch.Subcommands[name] = sub
	return
}

func (s *SlashCommand) AddSubcommandGroup(name, description string) (sub *SubSlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Subcommand name")
	}

	ch, ok := s.Handler.(*SubcommandHandler)
	if !ok {
		return nil, fmt.Errorf("Slash Command '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("Subcommand '%v' already exists", name)
	}

	sub = &SubSlashCommand{
		Name:        name,
		Description: description,
		IsGroup:     true,
		Handler:     &SubcommandHandler{Subcommands: make(map[string]*SubSlashCommand)},
	}

	ch.Subcommands[name] = sub
	return
}

type SubSlashCommand struct {
	Name        string
	Description string
	IsGroup     bool
	Handler     CommandHandler
}

func (s *SubSlashCommand) AddSubcommand(name, description string, handler func(h *Harmonia, i *Invocation)) (sub *SubSlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Subcommand name")
	}

	ch, ok := s.Handler.(*SubcommandHandler)
	if !ok {
		return nil, fmt.Errorf("Subcommand '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("Subcommand '%v' already exists", name)
	}

	sub = &SubSlashCommand{
		Name:        name,
		Description: description,
		Handler:     &SingleCommandHandler{Handler: handler},
	}

	ch.Subcommands[name] = sub
	return
}

type Invocation struct {
	*discordgo.Interaction
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *Author
	options []*discordgo.ApplicationCommandInteractionDataOption
}

func (i *Invocation) GetOptionMap() map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.options))
	for _, opt := range i.options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func (i *Invocation) GetOption(name string) *discordgo.ApplicationCommandInteractionDataOption {
	option, ok := i.GetOptionMap()[name]
	if !ok {
		return nil
	}
	return option
}

type Followup struct {
	*discordgo.Message
	Interaction *discordgo.Interaction
	Channel     *discordgo.Channel
	Guild       *discordgo.Guild
}

type Author struct {
	*discordgo.User
	IsMember     bool
	Guild        *discordgo.Guild
	JoinedAt     time.Time
	Nick         string
	Deaf         bool
	Mute         bool
	Roles        []*discordgo.Role
	PremiumSince *time.Time
}
