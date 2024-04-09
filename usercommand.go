package harmonia

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// A UserCommand describes a user command or USER application command.
type UserCommand struct {
	name               string
	guildID            string
	dmPermission       bool
	defaultPermissions int64

	commandFunc CommandFunc

	registration *discordgo.ApplicationCommand
}

// NewUserCommand returns a UserCommand with a given name
func NewUserCommand(name string) *UserCommand {
	if name == "" {
		log.Panic("empty command name")
	}

	return &UserCommand{
		name: name,
	}
}

// WithGuildID changes the guildID of the UserCommand and returns itself, so that it can be chained.
func (s *UserCommand) WithGuildID(guildID string) *UserCommand {
	s.guildID = guildID
	return s
}

// WithDMPermission changes the DM Permission of the UserCommand and returns itself, so that it can be chained.
func (s *UserCommand) WithDMPermission(isAllowed bool) *UserCommand {
	s.dmPermission = isAllowed
	return s
}

// WithDefaultPermissions changes the DefaultPermissions of the UserCommand and returns itself, so that it can be chained.
func (s *UserCommand) WithDefaultPermissions(defaultPermissions int64) *UserCommand {
	s.defaultPermissions = defaultPermissions
	return s
}

// WithCommand changes the CommandFunc that is called when the UserCommand is executed and returns itself, so that it can be chained.
func (s *UserCommand) WithCommand(commandFunc CommandFunc) *UserCommand {
	s.commandFunc = commandFunc
	return s
}

func (s *UserCommand) GetName() string {
	return s.name
}

func (s *UserCommand) Do(h *Harmonia, i *Invocation) {
	go s.commandFunc(h, i)
}

func (s *UserCommand) getRegistration() *discordgo.ApplicationCommand {
	if s.registration != nil {
		return s.registration
	}

	return &discordgo.ApplicationCommand{
		Name:    s.name,
		GuildID: s.guildID,
		Type:    discordgo.UserApplicationCommand,
	}
}

func (s *UserCommand) setRegistration(registration *discordgo.ApplicationCommand) {
	s.registration = registration
}
