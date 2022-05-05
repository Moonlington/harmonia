package harmonia

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Harmonia struct {
	*discordgo.Session
	Commands map[string]*SlashCommand
	running  bool
}

func New(token string) (h *Harmonia, err error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	h = &Harmonia{
		Session:  s,
		Commands: make(map[string]*SlashCommand),
	}

	return h, err
}

func (h *Harmonia) AddSlashCommand(name, description string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Slash Command name")
	}

	if _, ok := h.Commands[name]; ok {
		return nil, fmt.Errorf("Slash Command '%v' already exists", name)
	}

	c = &SlashCommand{
		Name:        name,
		Description: description,
		GuildID:     "",
		Handler:     &SingleCommandHandler{Handler: handler},
	}

	h.Commands[name] = c
	return
}

func (h *Harmonia) AddSlashCommandInGuild(name, description, GuildID string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	c, err = h.AddSlashCommand(name, description, handler)
	c.GuildID = GuildID
	return
}

func (h *Harmonia) AddSlashCommandWithSubcommands(name, description string) (c *SlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Slash Command name")
	}

	if _, ok := h.Commands[name]; ok {
		return nil, fmt.Errorf("Slash Command '%v' already exists", name)
	}

	c = &SlashCommand{
		Name:        name,
		Description: description,
		GuildID:     "",
		Handler:     &SubcommandHandler{Subcommands: make(map[string]*SubSlashCommand)},
	}

	h.Commands[name] = c
	return
}

func (h *Harmonia) AddSlashCommandWithSubcommandsInGuild(name, description, GuildID string) (c *SlashCommand, err error) {
	c, err = h.AddSlashCommandWithSubcommands(name, description)
	c.GuildID = GuildID
	return
}

func (h *Harmonia) authorFromInteraction(i *discordgo.Interaction) (a *Author, err error) {
	if i.Member == nil {
		return &Author{User: i.User, IsMember: false}, nil
	}

	// TODO: Error checking
	guild, _ := h.Guild(i.Member.GuildID)

	guildroles, _ := h.GuildRoles(i.Member.GuildID)
	roles := make([]*discordgo.Role, len(i.Member.Roles))

	j := 0
	for _, r := range guildroles {
		for _, mr := range i.Member.Roles {
			if r.ID == mr {
				roles[j] = r
				j++
				break
			}
		}
	}

	a = &Author{User: i.Member.User,
		IsMember:     true,
		Guild:        guild,
		JoinedAt:     i.Member.JoinedAt,
		Nick:         i.Member.Nick,
		Deaf:         i.Member.Deaf,
		Mute:         i.Member.Mute,
		Roles:        roles,
		PremiumSince: i.Member.PremiumSince}
	a.Avatar = i.Member.Avatar
	return a, nil
}

func (h *Harmonia) Respond(i *Invocation, content string) error {
	return h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func (h *Harmonia) Run() error {
	h.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if sc, ok := h.Commands[i.ApplicationCommandData().Name]; ok {
			//TODO: Error checking for each AND work out some way to use State in this.
			guild, _ := h.Guild(i.GuildID)
			channel, _ := h.Channel(i.ChannelID)
			author, _ := h.authorFromInteraction(i.Interaction)
			options := i.ApplicationCommandData().Options

			sc.Handler.Do(h, &Invocation{
				Interaction: i.Interaction,
				Guild:       guild,
				Channel:     channel,
				Author:      author,
				Options:     options,
			})
		}
	})

	err := h.Open()
	if err != nil {
		return err
	}

	for _, command := range h.Commands {
		cmd, err := h.ApplicationCommandCreate(h.State.User.ID, command.GuildID, &discordgo.ApplicationCommand{
			Name:        command.Name,
			Description: command.Description,
			Options:     command.Handler.GetOptions(),
		})
		if err != nil {
			return err
		}
		command.registration = cmd
	}
	return nil
}

func (h *Harmonia) RemoveCommand(name string) error {
	command, ok := h.Commands[name]
	if !ok {
		return fmt.Errorf("Command '%v' was not found", name)
	}

	if command.registration == nil {
		return fmt.Errorf("Command '%v' was not registered", name)
	}

	err := h.ApplicationCommandDelete(h.State.User.ID, command.GuildID, command.registration.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *Harmonia) RemoveAllCommands() error {
	for _, v := range h.Commands {
		if v.registration == nil {
			continue
		}

		err := h.ApplicationCommandDelete(h.State.User.ID, v.GuildID, v.registration.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
