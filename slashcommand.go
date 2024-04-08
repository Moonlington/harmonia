package harmonia

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// A SlashCommand describes a slash command or CHAT_INPUT application command.
type SlashCommand struct {
	Name         string
	Description  string
	GuildID      string
	Handler      CommandHandler
	registration *discordgo.ApplicationCommand
}

// AddSubcommand adds a Subcommand to a SlashCommand, this can only be done if the SlashCommand was created as a Subcommand Group. See AddSlashCommandWithSubcommands for more information.
func (s *SlashCommand) AddSubcommand(name, description string, handler func(h *Harmonia, i *Invocation)) (sub *SlashSubcommand, err error) {
	if name == "" {
		return nil, errors.New("empty Subcommand name")
	}

	ch, ok := s.Handler.(*CommandGroupHandler)
	if !ok {
		return nil, fmt.Errorf("command '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("subcommand '%v' already exists", name)
	}

	sub = &SlashSubcommand{
		Name:        name,
		Description: description,
		Handler:     &SingleCommandHandler{Handler: handler},
	}

	ch.Subcommands[name] = sub
	return
}

// AddSubcommand adds a Subcommand Group to a SlashCommand.
func (s *SlashCommand) AddSubcommandGroup(name, description string) (sub *SlashSubcommand, err error) {
	if name == "" {
		return nil, errors.New("empty Subcommand name")
	}

	ch, ok := s.Handler.(*CommandGroupHandler)
	if !ok {
		return nil, fmt.Errorf("command '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("subcommand '%v' already exists", name)
	}

	sub = &SlashSubcommand{
		Name:        name,
		Description: description,
		IsGroup:     true,
		Handler:     &CommandGroupHandler{Subcommands: make(map[string]*SlashSubcommand)},
	}

	ch.Subcommands[name] = sub
	return
}

// SlashSubcommand describes a Subcommand or a Subcommand group to a SlashCommand.
type SlashSubcommand struct {
	Name        string
	Description string
	IsGroup     bool
	Handler     CommandHandler
}

// AddSubcommand adds a Subcommand to a SubSlashCommand.
func (s *SlashSubcommand) AddSubcommand(name, description string, handler func(h *Harmonia, i *Invocation)) (sub *SlashSubcommand, err error) {
	if name == "" {
		return nil, errors.New("empty Subcommand name")
	}

	ch, ok := s.Handler.(*CommandGroupHandler)
	if !ok {
		return nil, fmt.Errorf("subcommand '%v' does not have a SubcommandHandler", s.Name)
	}

	if _, ok := ch.Subcommands[name]; ok {
		return nil, fmt.Errorf("subcommand '%v' already exists", name)
	}

	sub = &SlashSubcommand{
		Name:        name,
		Description: description,
		Handler:     &SingleCommandHandler{Handler: handler},
	}

	ch.Subcommands[name] = sub
	return
}

// An Invocation describes an incoming Interaction.
type Invocation struct {
	*discordgo.Interaction
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *Author

	options []*discordgo.ApplicationCommandInteractionDataOption

	// Only when the incoming Interaction is from a SelectMenu component.
	Values []string
}

// GetOptionMap returns a map of options passed through the Invocation.
func (i *Invocation) GetOptionMap() map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.options))
	for _, opt := range i.options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

// GetOption returns a specific option from an Invocation.
func (i *Invocation) GetOption(name string) *discordgo.ApplicationCommandInteractionDataOption {
	option, ok := i.GetOptionMap()[name]
	if !ok {
		return nil
	}
	return option
}
