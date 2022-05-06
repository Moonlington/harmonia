package harmonia

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// VERSION of Harmonia, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.2.0"

type Harmonia struct {
	*discordgo.Session
	Commands          map[string]*SlashCommand
	ComponentHandlers map[string]func(h *Harmonia, i *Invocation)
	running           bool
}

func New(token string) (h *Harmonia, err error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	h = &Harmonia{
		Session:           s,
		Commands:          make(map[string]*SlashCommand),
		ComponentHandlers: make(map[string]func(h *Harmonia, i *Invocation)),
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

func (h *Harmonia) interactionMessageFromMessage(m *discordgo.Message, i *discordgo.Interaction) *InteractionMessage {
	f := &InteractionMessage{Message: m, Interaction: i}

	if m != nil {
		guild, _ := h.Guild(m.GuildID)
		f.Guild = guild

		channel, _ := h.Channel(m.ChannelID)
		f.Channel = channel
	}

	return f
}

func (h *Harmonia) Respond(i *Invocation, content string) (*InteractionMessage, error) {
	err := h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) EphemeralRespond(i *Invocation, content string) (*InteractionMessage, error) {
	err := h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	})
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) RespondWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	err := h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: comp,
		},
	})
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) RespondComplex(i *Invocation, resp *discordgo.InteractionResponse) (*InteractionMessage, error) {
	err := h.InteractionRespond(i.Interaction, resp)
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) DeferRespond(i *Invocation) error {
	return h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func (h *Harmonia) EditResponse(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: content,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) EditResponseWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    content,
		Components: comp,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) DeleteResponse(i *Invocation) error {
	return h.InteractionResponseDelete(i.Interaction)
}

func (h *Harmonia) Followup(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) EphemeralFollowup(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Flags:   1 << 6,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) FollowupWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content:    content,
		Components: comp,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) FollowupComplex(i *Invocation, params *discordgo.WebhookParams) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, params)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

func (h *Harmonia) EditFollowup(f *InteractionMessage, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageEdit(f.Interaction, f.ID, &discordgo.WebhookEdit{
		Content: content,
	})
	return h.interactionMessageFromMessage(m, f.Interaction), err
}

func (h *Harmonia) EditFollowupWithComponents(f *InteractionMessage, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.FollowupMessageEdit(f.Interaction, f.ID, &discordgo.WebhookEdit{
		Content:    content,
		Components: comp,
	})
	return h.interactionMessageFromMessage(m, f.Interaction), err
}

func (h *Harmonia) DeleteFollowup(f *InteractionMessage) error {
	return h.FollowupMessageDelete(f.Interaction, f.ID)
}

func (h *Harmonia) AddComponentHandler(customID string, handler func(h *Harmonia, i *Invocation)) error {
	if customID == "" {
		return errors.New("Empty CustomID")
	}

	if _, ok := h.ComponentHandlers[customID]; ok {
		return fmt.Errorf("CustomID '%v' already exists", customID)
	}

	h.ComponentHandlers[customID] = handler
	return nil
}

func (h *Harmonia) AddComponentHandlerToInteractionMessage(f *InteractionMessage, customID string, handler func(h *Harmonia, i *Invocation)) error {
	if customID == "" {
		return errors.New("Empty CustomID")
	}

	followupcustomID := fmt.Sprintf("%v-%v", f.ID, customID)

	if _, ok := h.ComponentHandlers[followupcustomID]; ok {
		return fmt.Errorf("CustomID '%v' already exists on Followup '%v'", customID, f.ID)
	}

	h.ComponentHandlers[followupcustomID] = handler
	return nil
}

func (h *Harmonia) RemoveComponentHandler(customID string) error {
	if _, ok := h.ComponentHandlers[customID]; !ok {
		return fmt.Errorf("CustomID '%v' not found", customID)
	}
	delete(h.ComponentHandlers, customID)
	return nil
}

func (h *Harmonia) RemoveComponentHandlerFromInteractionMessage(f *InteractionMessage, customID string) error {
	followupcustomID := fmt.Sprintf("%v-%v", f.ID, customID)
	if _, ok := h.ComponentHandlers[followupcustomID]; !ok {
		return fmt.Errorf("CustomID '%v' not found on Followup '%v'", customID, f.ID)
	}
	delete(h.ComponentHandlers, followupcustomID)
	return nil
}

func (h *Harmonia) Run() error {
	h.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if slashCommand, ok := h.Commands[i.ApplicationCommandData().Name]; ok {
				//TODO: Error checking for each AND work out some way to use State in this.
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := h.authorFromInteraction(i.Interaction)
				options := i.ApplicationCommandData().Options

				slashCommand.Handler.Do(h, &Invocation{
					Interaction: i.Interaction,
					Guild:       guild,
					Channel:     channel,
					Author:      author,
					options:     options,
				})
			}
			return
		case discordgo.InteractionMessageComponent:
			if componentHandler, ok := h.ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := h.authorFromInteraction(i.Interaction)
				values := i.MessageComponentData().Values

				componentHandler(h, &Invocation{
					Interaction: i.Interaction,
					Guild:       guild,
					Channel:     channel,
					Author:      author,
					Values:      values,
				})
				return
			}

			followupcustomID := fmt.Sprintf("%v-%v", i.Message.ID, i.MessageComponentData().CustomID)

			if componentHandler, ok := h.ComponentHandlers[followupcustomID]; ok {
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := h.authorFromInteraction(i.Interaction)
				values := i.MessageComponentData().Values

				componentHandler(h, &Invocation{
					Interaction: i.Interaction,
					Guild:       guild,
					Channel:     channel,
					Author:      author,
					Values:      values,
				})
				return
			}
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

	delete(h.Commands, name)
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
