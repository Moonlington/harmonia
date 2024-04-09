package harmonia

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// GroupSlashCommand describes a group of slash commands.
type GroupSlashCommand struct {
	name               string
	description        string
	guildID            string
	dmPermission       bool
	defaultPermissions int64

	subcommands map[string]CommandHandler

	registration *discordgo.ApplicationCommand
}

// NewGroupSlashCommand returns a GroupSlashCommand with a given name
func NewGroupSlashCommand(name string) *GroupSlashCommand {
	if name == "" {
		log.Panic("empty command name")
	}

	if !slashCommandNameRegex.MatchString(name) {
		log.Panic("slash command name does not match with the CHAT_INPUT regex.")
	}

	return &GroupSlashCommand{
		name:        name,
		subcommands: make(map[string]CommandHandler),
	}
}

// WithDescription changes the description of the GroupSlashCommand and returns itself, so that it can be chained.
func (s *GroupSlashCommand) WithDescription(description string) *GroupSlashCommand {
	s.description = description
	return s
}

// WithGuildID changes the guildID of the GroupSlashCommand and returns itself, so that it can be chained.
func (s *GroupSlashCommand) WithGuildID(guildID string) *GroupSlashCommand {
	s.guildID = guildID
	return s
}

func (s *GroupSlashCommand) WithSubCommands(subcommands ...CommandHandler) *GroupSlashCommand {
	for _, command := range subcommands {
		name := command.GetName()

		if _, ok := command.(*SlashCommand); !ok {
			if _, ok := command.(*GroupSlashCommand); !ok {
				log.Panic("supplied subcommand is neither SlashCommand nor GroupSlashCommand")
			}
		}

		if _, ok := s.subcommands[name]; ok {
			log.Panic("duplicate subcommand name")
		}

		s.subcommands[name] = command
	}
	return s
}

// WithDMPermission changes the DM Permission of the GroupSlashCommand and returns itself, so that it can be chained.
func (s *GroupSlashCommand) WithDMPermission(isAllowed bool) *GroupSlashCommand {
	s.dmPermission = isAllowed
	return s
}

// WithDefaultPermissions changes the DefaultPermissions of the GroupSlashCommand and returns itself, so that it can be chained.
func (s *GroupSlashCommand) WithDefaultPermissions(defaultPermissions int64) *GroupSlashCommand {
	s.defaultPermissions = defaultPermissions
	return s
}

func (s *GroupSlashCommand) GetName() string {
	return s.name
}

func (s *GroupSlashCommand) Do(h *Harmonia, i *Invocation) {
	options := i.options
	if command, ok := s.subcommands[options[0].Name]; ok {
		i.options = options[0].Options
		go command.Do(h, i)
	}
}

func (s *GroupSlashCommand) getRegistration() *discordgo.ApplicationCommand {
	if s.registration != nil {
		return s.registration
	}

	options := make([]*discordgo.ApplicationCommandOption, len(s.subcommands))
	i := 0
	for _, command := range s.subcommands {
		t := discordgo.ApplicationCommandOptionSubCommand
		if _, ok := command.(*GroupSlashCommand); ok {
			t = discordgo.ApplicationCommandOptionSubCommandGroup
		}

		data := command.getRegistration()

		options[i] = &discordgo.ApplicationCommandOption{
			Name:        data.Name,
			Description: data.Description,
			Options:     data.Options,
			Type:        t,
		}
		i++
	}

	return &discordgo.ApplicationCommand{
		Name:                     s.name,
		Description:              s.description,
		GuildID:                  s.guildID,
		Options:                  options,
		DMPermission:             &s.dmPermission,
		DefaultMemberPermissions: &s.defaultPermissions,
		Type:                     discordgo.ChatApplicationCommand,
	}
}

func (s *GroupSlashCommand) setRegistration(registration *discordgo.ApplicationCommand) {
	s.registration = registration
}
