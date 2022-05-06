package harmonia

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// A SlashCommand describes a slash command or application command.
type SlashCommand struct {
	Name         string
	Description  string
	GuildID      string
	Handler      CommandHandler
	registration *discordgo.ApplicationCommand
}

// AddSubcommand adds a Subcommand to a SlashCommand, this can only be done if the SlashCommand was created as a Subcommand Group. See AddSlashCommandWithSubcommands for more information.
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

// AddSubcommand adds a Subcommand Group to a SlashCommand.
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

// SubSlashCommand describes a Subcommand or a Subcommand group to a SlashCommand.
type SubSlashCommand struct {
	Name        string
	Description string
	IsGroup     bool
	Handler     CommandHandler
}

// AddSubcommand adds a Subcommand to a SubSlashCommand.
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

// InteractionMessage describes a message sent as a follow-up or response to an Interaction.
type InteractionMessage struct {
	*discordgo.Message
	Interaction *discordgo.Interaction
	Channel     *discordgo.Channel
	Guild       *discordgo.Guild
}

// An Author describes either a User or Member, depending if the message was sent in a Guild or DMs.
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
