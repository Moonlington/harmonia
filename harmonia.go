package harmonia

import (
	"errors"

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
		Commands: map[string]*SlashCommand{},
	}

	return h, err
}

func (h *Harmonia) AddSlashCommand(name, description string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Slash Command name")
	}

	if _, ok := h.Commands[name]; ok {
		return nil, errors.New("Duplicate Slash Command name")
	}

	c = &SlashCommand{
		Name:        name,
		Description: description,
		Handler:     handler,
		Options:     []*discordgo.ApplicationCommandOption{},
		GuildID:     "",
		// Subcommands: map[string]*SlashCommand{},
	}

	h.Commands[name] = c
	return
}

func (h *Harmonia) AddSlashCommandInGuild(name, description, GuildID string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	c, err = h.AddSlashCommand(name, description, handler)
	c.GuildID = GuildID
	return
}

func (h *Harmonia) AuthorFromInteraction(i *discordgo.Interaction) (a *Author, err error) {
	if i.Member == nil {
		return &Author{User: i.User, IsMember: false}, nil
	}
	guild, _ := h.State.Guild(i.Member.GuildID)
	roles := make([]*discordgo.Role, len(i.Member.Roles))
	for j, r := range i.Member.Roles {
		role, _ := h.State.Role(i.Member.GuildID, r)
		roles[j] = role
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
			guild, _ := h.State.Guild(i.GuildID)
			channel, _ := h.State.Channel(i.ChannelID)
			author, _ := h.AuthorFromInteraction(i.Interaction)
			sc.Handler(h, &Invocation{
				Interaction: i.Interaction,
				Guild:       guild,
				Channel:     channel,
				Author:      author,
			})
		}
	})

	err := h.Open()
	if err != nil {
		return err
	}

	for _, v := range h.Commands {
		cmd, err := h.ApplicationCommandCreate(h.State.User.ID, v.GuildID, &discordgo.ApplicationCommand{
			Name:        v.Name,
			Description: v.Description,
			Options:     v.Options,
		})
		if err != nil {
			return err
		}
		v.registration = cmd
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
