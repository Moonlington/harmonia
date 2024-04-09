package harmonia

import "github.com/bwmarrin/discordgo"

// An Invocation describes an incoming Interaction.
type Invocation struct {
	*discordgo.Interaction
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *Author

	options []*discordgo.ApplicationCommandInteractionDataOption

	// Only when the incoming Interaction is from a SelectMenu component.
	Values []string
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
