package harmonia

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// A MessageCommand describes a user command or USER application command.
type MessageCommand struct {
	name               string
	guildID            string
	dmPermission       bool
	defaultPermissions *int64

	commandFunc CommandFunc

	registration *discordgo.ApplicationCommand
}

// NewMessageCommand returns a MessageCommand with a given name
func NewMessageCommand(name string) *MessageCommand {
	if name == "" {
		log.Panic("empty command name")
	}

	return &MessageCommand{
		name: name,
	}
}

// WithGuildID changes the guildID of the MessageCommand and returns itself, so that it can be chained.
func (s *MessageCommand) WithGuildID(guildID string) *MessageCommand {
	s.guildID = guildID
	return s
}

// WithDMPermission changes the DM Permission of the MessageCommand and returns itself, so that it can be chained.
func (s *MessageCommand) WithDMPermission(isAllowed bool) *MessageCommand {
	s.dmPermission = isAllowed
	return s
}

// WithDefaultPermissions changes the DefaultPermissions of the MessageCommand and returns itself, so that it can be chained.
func (s *MessageCommand) WithDefaultPermissions(defaultPermissions int64) *MessageCommand {
	s.defaultPermissions = &defaultPermissions
	return s
}

// WithCommand changes the CommandFunc that is called when the MessageCommand is executed and returns itself, so that it can be chained.
func (s *MessageCommand) WithCommand(commandFunc CommandFunc) *MessageCommand {
	s.commandFunc = commandFunc
	return s
}

func (s *MessageCommand) GetName() string {
	return s.name
}

func (s *MessageCommand) Do(h *Harmonia, i *Invocation) {
	go s.commandFunc(h, i)
}

func (s *MessageCommand) getRegistration() *discordgo.ApplicationCommand {
	if s.registration != nil {
		return s.registration
	}

	return &discordgo.ApplicationCommand{
		Name:                     s.name,
		GuildID:                  s.guildID,
		DMPermission:             &s.dmPermission,
		DefaultMemberPermissions: s.defaultPermissions,
		Type:                     discordgo.MessageApplicationCommand,
	}
}

func (s *MessageCommand) setRegistration(registration *discordgo.ApplicationCommand) {
	s.registration = registration
}
