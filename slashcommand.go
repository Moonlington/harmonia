package harmonia

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var slashCommandNameRegex = regexp.MustCompile(`^[-_\p{L}\p{N}]{1,32}$`)

// A SlashCommand describes a slash command or CHAT_INPUT application command.
type SlashCommand struct {
	name               string
	description        string
	guildID            string
	dmPermission       bool
	defaultPermissions *int64

	commandFunc CommandFunc
	options     []*Option

	registration *discordgo.ApplicationCommand
}

// NewSlashCommand returns a SlashCommand with a given name
func NewSlashCommand(name string) *SlashCommand {
	if name == "" {
		log.Panic("empty command name")
	}

	if !slashCommandNameRegex.MatchString(name) {
		log.Panic("slash command name does not match with the CHAT_INPUT regex.")
	}

	return &SlashCommand{
		name: name,
	}
}

// WithDescription changes the description of the SlashCommand and returns itself, so that it can be chained.
func (s *SlashCommand) WithDescription(description string) *SlashCommand {
	s.description = description
	return s
}

// WithGuildID changes the guildID of the SlashCommand and returns itself, so that it can be chained.
func (s *SlashCommand) WithGuildID(guildID string) *SlashCommand {
	s.guildID = guildID
	return s
}

// WithDMPermission changes the DM Permission of the SlashCommand and returns itself, so that it can be chained.
func (s *SlashCommand) WithDMPermission(isAllowed bool) *SlashCommand {
	s.dmPermission = isAllowed
	return s
}

// WithDefaultPermissions changes the DefaultPermissions of the SlashCommand and returns itself, so that it can be chained.
func (s *SlashCommand) WithDefaultPermissions(defaultPermissions int64) *SlashCommand {
	s.defaultPermissions = &defaultPermissions
	return s
}

// WithCommand changes the CommandFunc that is called when the SlashCommand is executed and returns itself, so that it can be chained.
func (s *SlashCommand) WithCommand(commandFunc CommandFunc) *SlashCommand {
	s.commandFunc = commandFunc
	return s
}

// WithOptions changes the options in the SlashCommand and returns itself, so that it can be chained.
func (s *SlashCommand) WithOptions(options ...*Option) *SlashCommand {
	s.options = options
	return s
}

func (s *SlashCommand) GetName() string {
	return s.name
}

func (s *SlashCommand) Do(h *Harmonia, i *Invocation) {
	go s.commandFunc(h, i)
}

func (s *SlashCommand) getRegistration() *discordgo.ApplicationCommand {
	if s.registration != nil {
		return s.registration
	}

	options := make([]*discordgo.ApplicationCommandOption, len(s.options))
	for i, v := range s.options {
		options[i] = v.ApplicationCommandOption
	}

	return &discordgo.ApplicationCommand{
		Name:                     s.name,
		Description:              s.description,
		GuildID:                  s.guildID,
		Options:                  options,
		DMPermission:             &s.dmPermission,
		DefaultMemberPermissions: s.defaultPermissions,
		Type:                     discordgo.ChatApplicationCommand,
	}
}

func (s *SlashCommand) setRegistration(registration *discordgo.ApplicationCommand) {
	s.registration = registration
}
