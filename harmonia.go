package harmonia

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// VERSION of Harmonia, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.7.0"

// A Harmonia represents a connection to the Discord API and contains the slash commands and component handlers used by Harmonia.
type Harmonia struct {
	*discordgo.Session
	Commands          map[string]CommandHandler
	ComponentHandlers map[string]CommandFunc
}

// New creates a new Discord session with the provided token and wraps the Harmonia struct around it.
func New(token string) (h *Harmonia, err error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	h = &Harmonia{
		Session:           s,
		Commands:          make(map[string]CommandHandler),
		ComponentHandlers: make(map[string]CommandFunc),
	}

	return h, err
}

// AddCommand adds a command to Harmonia.
func (h *Harmonia) AddCommand(command CommandHandler) (err error) {
	name := command.GetName()
	if _, ok := h.Commands[name]; ok {
		return fmt.Errorf("command '%v' already exists", name)
	}

	h.Commands[name] = command
	return
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

// RespondComplex allows you full freedom to respond with whatever you'd like.
func (h *Harmonia) RespondComplex(i *Invocation, resp *discordgo.InteractionResponse) (*InteractionMessage, error) {
	err := h.InteractionRespond(i.Interaction, resp)
	if err != nil {
		return nil, err
	}
	m, err := h.InteractionResponse(i.Interaction)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// Respond allows Harmonia to easily respond to an Invocation with a string.
func (h *Harmonia) Respond(i *Invocation, content string) (*InteractionMessage, error) {
	return h.RespondComplex(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

// EphemeralRespond does the same as Respond, but sets the flag such that only the invoker can see the message.
func (h *Harmonia) EphemeralRespond(i *Invocation, content string) (*InteractionMessage, error) {
	return h.RespondComplex(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	})
}

// RespondWithComponents does the same as Respond, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the response.
func (h *Harmonia) RespondWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	return h.RespondComplex(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: comp,
		},
	})
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

// FollowupComplex allows you full freedom to follow-up with whatever you'd like.
func (h *Harmonia) FollowupComplex(i *Invocation, params *discordgo.WebhookParams) (*InteractionMessage, error) {
	m, err := h.FollowupMessageCreate(i.Interaction, true, params)
	return h.interactionMessageFromMessage(m, i.Interaction), err
}

// Followup sends a follow-up message to the Interaction, this does require you to have used DeferResponse before.
func (h *Harmonia) Followup(i *Invocation, content string) (*InteractionMessage, error) {
	return h.FollowupComplex(i, &discordgo.WebhookParams{
		Content: content,
	})
}

// EphemeralFollowup does the same as Followup, but sets the flag such that only the invoker can see the message.
func (h *Harmonia) EphemeralFollowup(i *Invocation, content string) (*InteractionMessage, error) {
	return h.FollowupComplex(i, &discordgo.WebhookParams{
		Content: content,
		Flags:   1 << 6,
	})
}

// FollowupWithComponents does the same as Followup, but also takes in a 2D slice of discordgo.MessageComponents that will be added to the response.
func (h *Harmonia) FollowupWithComponents(i *Invocation, content string, components [][]discordgo.MessageComponent) (*InteractionMessage, error) {
	comp := ParseComponentMatrix(components)
	return h.FollowupComplex(i, &discordgo.WebhookParams{
		Content:    content,
		Components: comp,
	})
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
func (h *Harmonia) AddComponentHandler(customID string, handler CommandFunc) error {
	if customID == "" {
		return errors.New("empty CustomID")
	}

	if _, ok := h.ComponentHandlers[customID]; ok {
		return fmt.Errorf("CustomID '%v' already exists", customID)
	}

	h.ComponentHandlers[customID] = handler
	return nil
}

// AddComponentHandlerToInteractionMessage adds a handler for a component, but will be handled only on its original Interaction.
// This is done by prepending the InteractionMessage's ID to the customID. Harmonia will do the heavy lifting from there.
func (h *Harmonia) AddComponentHandlerToInteractionMessage(f *InteractionMessage, customID string, handler CommandFunc) error {
	if customID == "" {
		return errors.New("empty CustomID")
	}

	followupcustomID := fmt.Sprintf("%v-%v", f.ID, customID)

	if _, ok := h.ComponentHandlers[followupcustomID]; ok {
		return fmt.Errorf("customID '%v' already exists on Followup '%v'", customID, f.ID)
	}

	h.ComponentHandlers[followupcustomID] = handler
	return nil
}

// RemoveComponentHandler removes a component handler.
func (h *Harmonia) RemoveComponentHandler(customID string) error {
	if _, ok := h.ComponentHandlers[customID]; !ok {
		return fmt.Errorf("customID '%v' not found", customID)
	}
	delete(h.ComponentHandlers, customID)
	return nil
}

// RemoveComponentHandlerFromInteractionMessage removes a component handler from an InteractionMessage.
func (h *Harmonia) RemoveComponentHandlerFromInteractionMessage(f *InteractionMessage, customID string) error {
	followupcustomID := fmt.Sprintf("%v-%v", f.ID, customID)
	if _, ok := h.ComponentHandlers[followupcustomID]; !ok {
		return fmt.Errorf("customID '%v' not found on Followup '%v'", customID, f.ID)
	}
	delete(h.ComponentHandlers, followupcustomID)
	return nil
}

// Run starts the Harmonia bot up and does the handling for slash commands and components for you.
func (h *Harmonia) Run() error {
	h.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if command, ok := h.Commands[i.ApplicationCommandData().Name]; ok {
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := AuthorFromInteraction(h, i.Interaction)
				options := i.ApplicationCommandData().Options

				command.Do(h, &Invocation{
					Interaction: i.Interaction,
					Guild:       guild,
					Channel:     channel,
					Author:      author,
					options:     options,
					targetID:    i.ApplicationCommandData().TargetID,
				})
			}
			return
		case discordgo.InteractionMessageComponent:
			if componentHandler, ok := h.ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				guild, _ := h.Guild(i.GuildID)
				channel, _ := h.Channel(i.ChannelID)
				author, _ := AuthorFromInteraction(h, i.Interaction)
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
				author, _ := AuthorFromInteraction(h, i.Interaction)
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
		data := command.getRegistration()
		registration, err := h.ApplicationCommandCreate(h.State.User.ID, data.GuildID, data)
		if err != nil {
			return err
		}
		command.setRegistration(registration)
	}
	return nil
}

// RemoveCommand removes a slash command from Harmonia and from the Discord API.
func (h *Harmonia) RemoveCommand(name string) error {
	command, ok := h.Commands[name]
	if !ok {
		return fmt.Errorf("command '%v' was not found", name)
	}

	registration := command.getRegistration()

	if registration.ID == "" {
		return fmt.Errorf("command '%v' was not registered", name)
	}

	err := h.ApplicationCommandDelete(h.State.User.ID, registration.GuildID, registration.ID)
	if err != nil {
		return err
	}

	delete(h.Commands, name)
	return nil
}

// RemoveAllCommands does removes all registered commands from the Discord API.
func (h *Harmonia) RemoveAllCommands() error {
	globals, err := h.ApplicationCommands(h.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, global := range globals {
		err := h.ApplicationCommandDelete(h.State.User.ID, global.GuildID, global.ID)
		if err != nil {
			return err
		}
	}

	for _, guild := range h.State.Guilds {
		locals, err := h.ApplicationCommands(h.State.User.ID, guild.ID)
		if err != nil {
			return err
		}

		for _, local := range locals {
			err := h.ApplicationCommandDelete(h.State.User.ID, local.GuildID, local.ID)
			if err != nil {
				return err
			}
		}
	}
	h.Commands = make(map[string]CommandHandler)
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

// InteractionMessage describes a message sent as a follow-up or response to an Interaction.
type InteractionMessage struct {
	*discordgo.Message
	Interaction *discordgo.Interaction
	Channel     *discordgo.Channel
	Guild       *discordgo.Guild
}
