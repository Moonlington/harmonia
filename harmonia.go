package harmonia

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// VERSION of Harmonia, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.5.3"

var slashCommandNameRegex = regexp.MustCompile(`^[-_\p{L}\p{N}]{1,32}$`)

// A Harmonia represents a connection to the Discord API and contains the slash commands and component handlers used by Harmonia.
type Harmonia struct {
	*discordgo.Session
	Commands          map[string]*SlashCommand
	ComponentHandlers map[string]func(h *Harmonia, i *Invocation)
}

// New creates a new Discord session with the provided token and wraps the Harmonia struct around it.
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

// AddSlashCommand adds a slash command to Harmonia.
func (h *Harmonia) AddSlashCommand(name, description string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Slash Command name")
	}

	if !slashCommandNameRegex.MatchString(name) {
		return nil, errors.New("Slash Command name does not match with the CHAT_INPUT regex.")
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

// GuildAddSlashCommand does the same as AddSlashCommand, but only adds it for a specific GuildID.
func (h *Harmonia) GuildAddSlashCommand(name, description, GuildID string, handler func(h *Harmonia, i *Invocation)) (c *SlashCommand, err error) {
	c, err = h.AddSlashCommand(name, description, handler)
	c.GuildID = GuildID
	return
}

// AddSlashCommandWithSubcommands adds a subcommand group, it itself has no handler, but you can use the returned SlashCommand to add Subcommands to the SlashCommand.
func (h *Harmonia) AddSlashCommandWithSubcommands(name, description string) (c *SlashCommand, err error) {
	if name == "" {
		return nil, errors.New("Empty Slash Command name")
	}

	if !slashCommandNameRegex.MatchString(name) {
		return nil, errors.New("Slash Command name does not match with the CHAT_INPUT regex.")
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

// GuildAddSlashCommandWithSubcommands does the same as AddSlashCommandWithSubcommands, but only adds it for a specific GuildID.
func (h *Harmonia) GuildAddSlashCommandWithSubcommands(name, description, GuildID string) (c *SlashCommand, err error) {
	c, err = h.AddSlashCommandWithSubcommands(name, description)
	c.GuildID = GuildID
	return
}

// AuthorFromInteraction uses the information obtained from the Interaction to create an Author.
func (h *Harmonia) AuthorFromInteraction(i *discordgo.Interaction) (a *Author, err error) {
	if i.Member == nil {
		return AuthorFromUser(i.User), nil
	}

	i.Member.GuildID = i.GuildID
	return h.AuthorFromMember(i.Member)
}

// AuthorFromMember returns an Author from a *discordgo.Member.
func (h *Harmonia) AuthorFromMember(member *discordgo.Member) (*Author, error) {
	guild, err := h.Guild(member.GuildID)
	if err != nil {
		return nil, err
	}

	roles, err := h.RolesFromMember(member)
	if err != nil {
		return nil, err
	}

	a := &Author{User: member.User,
		IsMember:     true,
		Guild:        guild,
		JoinedAt:     member.JoinedAt,
		Nick:         member.Nick,
		Deaf:         member.Deaf,
		Mute:         member.Mute,
		Roles:        roles,
		PremiumSince: member.PremiumSince,
	}
	a.Avatar = member.Avatar
	return a, nil
}

// RolesFromMember returns a slice of *discordgo.Role from a *discordgo.Member.
func (h *Harmonia) RolesFromMember(member *discordgo.Member) ([]*discordgo.Role, error) {
	guildroles, err := h.GuildRoles(member.GuildID)
	if err != nil {
		return nil, err
	}

	roles := make([]*discordgo.Role, 0, len(member.Roles))
	for _, roleid := range member.Roles {
		for _, role := range guildroles {
			if role.ID == roleid {
				roles = append(roles, role)
			}
		}
	}

	return roles, nil
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

// Respond allows Harmonia to easily respond to an Invocation with a string.
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

// EphemeralRespond does the same as Respond, but sets the flag such that only the invoker can see the message.
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

// RespondWithComponents does the same as Respond, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the response.
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

// RespondComplex allows you full freedom to respond with whatever you'd like.
func (h *Harmonia) RespondComplex(i *Invocation, resp *discordgo.InteractionResponse) (*InteractionMessage, error) {
	err := h.InteractionRespond(i.Interaction, resp)
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// DeferResponse sends an acknowledgement to the DiscordAPI, allowing you to send a follow-up message later. See Followup for that.
func (h *Harmonia) DeferResponse(i *Invocation) error {
	return h.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// EditResponse edits an already sent response.
func (h *Harmonia) EditResponse(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// EditResponseWithComponents does the same as EditResponse, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the response.
func (h *Harmonia) EditResponseWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &comp,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// DeleteResponse deletes a response.
func (h *Harmonia) DeleteResponse(i *Invocation) error {
	return h.InteractionResponseDelete(i.Interaction)
}

// Followup sends a follow-up message to the Interaction, this does require you to have used DeferResponse before.
func (h *Harmonia) Followup(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// EphemeralFollowup does the same as Followup, but sets the flag such that only the invoker can see the message.
func (h *Harmonia) EphemeralFollowup(i *Invocation, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Flags:   1 << 6,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// FollowupWithComponents does the same as Followup, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the response.
func (h *Harmonia) FollowupWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content:    content,
		Components: comp,
	})
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// FollowupComplex allows you full freedom to follow-up with whatever you'd like.
func (h *Harmonia) FollowupComplex(i *Invocation, params *discordgo.WebhookParams) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, params)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// EditFollowup allows you to edit a follow-up message.
func (h *Harmonia) EditFollowup(f *InteractionMessage, content string) (*InteractionMessage, error) {
	m, err := h.FollowupMessageEdit(f.Interaction, f.ID, &discordgo.WebhookEdit{
		Content: &content,
	})
	return h.interactionMessageFromMessage(m, f.Interaction), err
}

// EditFollowupWithComponents does the same as EditFollowup, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the follow-up message.
func (h *Harmonia) EditFollowupWithComponents(f *InteractionMessage, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	m, err := h.FollowupMessageEdit(f.Interaction, f.ID, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &comp,
	})
	return h.interactionMessageFromMessage(m, f.Interaction), err
}

// DeleteFollowup deletes a follow-up message.
func (h *Harmonia) DeleteFollowup(f *InteractionMessage) error {
	return h.FollowupMessageDelete(f.Interaction, f.ID)
}

// AddComponentHandler adds a handler for a component.
// I suggest this is only used for globally used components, and not for components used on a message by message basis. See AddComponentHandlerToInteractionMessage
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

// AddComponentHandlerToInteractionMessage adds a handler for a component, but will be handled only on its original Interaction.
// This is done by prepending the InteractionMessage's ID to the customID. Harmonia will do the heavy lifting from there.
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

// RemoveComponentHandler removes a component handler.
func (h *Harmonia) RemoveComponentHandler(customID string) error {
	if _, ok := h.ComponentHandlers[customID]; !ok {
		return fmt.Errorf("CustomID '%v' not found", customID)
	}
	delete(h.ComponentHandlers, customID)
	return nil
}

// RemoveComponentHandlerFromInteractionMessage removes a component handler from an InteractionMessage.
func (h *Harmonia) RemoveComponentHandlerFromInteractionMessage(f *InteractionMessage, customID string) error {
	followupcustomID := fmt.Sprintf("%v-%v", f.ID, customID)
	if _, ok := h.ComponentHandlers[followupcustomID]; !ok {
		return fmt.Errorf("CustomID '%v' not found on Followup '%v'", customID, f.ID)
	}
	delete(h.ComponentHandlers, followupcustomID)
	return nil
}

// Run starts the Harmonia bot up and does the handling for slash commands and components for you.
func (h *Harmonia) Run() error {
	h.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if slashCommand, ok := h.Commands[i.ApplicationCommandData().Name]; ok {
				//TODO: Error checking for each AND work out some way to use State in this.
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := h.AuthorFromInteraction(i.Interaction)
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
				author, _ := h.AuthorFromInteraction(i.Interaction)
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
				author, _ := h.AuthorFromInteraction(i.Interaction)
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

// RemoveCommand removes a slash command from Harmonia and from the Discord API.
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

// RemoveAllCommands does removes all registered commands on this Harmonia instance and from the Discord API.
func (h *Harmonia) RemoveAllCommands() error {
	for _, command := range h.Commands {
		h.RemoveCommand(command.Name)
	}
	return nil
}

// ParseComponentMatrix parses a 2D slice of MessageComponents and returns a 1D slice of MessageComponents with ActionsRows.
func ParseComponentMatrix(components [][]discordgo.MessageComponent) []discordgo.MessageComponent {
	comp := make([]discordgo.MessageComponent, len(components))
	for i, c := range components {
		comp[i] = &discordgo.ActionsRow{Components: c}
	}
	return comp
}
