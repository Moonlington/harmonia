package harmonia

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// An Option is a wrapper around an ApplicationCommandOption with added functionality.
type Option struct {
	*discordgo.ApplicationCommandOption
}

// NewOption returns an option with given name and type.
func NewOption(name string, t discordgo.ApplicationCommandOptionType) *Option {
	if name == "" {
		log.Panic("empty option name")
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

// IsRequired sets the requirement of the Option to true and returns itself, so that it can be chained.
func (o *Option) IsRequired() *Option {
	o.Required = true
	return o
}

// AddChoice adds a choice to an option, value should be the same as the choice's type and returns itself, so that it can be chained.
func (o *Option) AddChoice(name string, value interface{}) *Option {
	c := &discordgo.ApplicationCommandOptionChoice{
		Name:  name,
		Value: value,
	}
	o.Choices = append(o.Choices, c)
	return o
}
