package harmonia

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// An Option is a wrapper around an ApplicationCommandOption with added functionality.
type Option struct {
	*discordgo.ApplicationCommandOption
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

// AddOption adds an Option to the SlashCommand.
func (s *SlashCommand) AddOption(name, description string, required bool, t discordgo.ApplicationCommandOptionType) (*Option, error) {
	if name == "" {
		return nil, errors.New("Empty Option name")
	}

	ch, ok := s.Handler.(*SingleCommandHandler)
	if !ok {
		return nil, fmt.Errorf("Slash Command '%v' is not a SingleCommandHandler and does not support Options", s.Name)
	}

	for _, v := range ch.Options {
		if v.Name == name {
			return nil, fmt.Errorf("Option '%v' already exists", name)
		}
	}

	o := &Option{&discordgo.ApplicationCommandOption{
		Type:        t,
		Name:        name,
		Description: description,
		Required:    required,
	}}

	ch.Options = append(ch.Options, o)
	return o, nil
}
