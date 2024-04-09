package harmonia

import (
	"github.com/bwmarrin/discordgo"
)

// An Option is a wrapper around an ApplicationCommandOption with added functionality.
type Option struct {
	*discordgo.ApplicationCommandOption
}

// NewOption returns an option with given name and type.
func NewOption(name string, t discordgo.ApplicationCommandOptionType) *Option {
	if name == "" {
		panic("empty option name")
	}

	return &Option{
		&discordgo.ApplicationCommandOption{
			Type: t,
			Name: name,
		},
	}
}

// WithDescription changes the description of the Option and returns itself, so that it can be chained.
func (o *Option) WithDescription(description string) *Option {
	o.Description = description
	return o
}

// WithRequired sets the requirement of the Option and returns itself, so that it can be chained.
func (o *Option) WithRequired(required bool) *Option {
	o.Required = required
	return o
}

// AddChoice adds a choice to an option, value should be the same as the choice's type.
func (o *Option) AddChoice(name string, value interface{}) *discordgo.ApplicationCommandOptionChoice {
	c := &discordgo.ApplicationCommandOptionChoice{
		Name:  name,
		Value: value,
	}
	o.Choices = append(o.Choices, c)
	return c
}
