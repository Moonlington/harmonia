package harmonia

import (
	"github.com/bwmarrin/discordgo"
)

// An Invocation describes an incoming Interaction.
type Invocation struct {
	*discordgo.Interaction
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *Author

	options []*discordgo.ApplicationCommandInteractionDataOption

	// Only when the incoming Interaction is from a SelectMenu component.
	Values []string

	// Only when the incoming Interaction is from a UserCommand or MessageCommand.
	targetID string
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

// TargetAuthor takes the targetID from the invocation and returns an Author struct from it.
func (i *Invocation) TargetAuthor(h *Harmonia) (*Author, error) {
	if i.Guild != nil {
		member, err := h.GuildMember(i.Guild.ID, i.targetID)
		if err != nil {
			return nil, err
		}

		return AuthorFromMember(h, member)
	}
	user, err := h.User(i.targetID)
	if err != nil {
		return nil, err
	}
	return AuthorFromUser(user), nil
}

// TargetMessage takes the targetID from the invocation and returns an Message struct from it.
func (i *Invocation) TargetMessage(h *Harmonia) (*discordgo.Message, error) {
	return h.ChannelMessage(i.ChannelID, i.targetID)
}
