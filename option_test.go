package harmonia

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestAddOption(t *testing.T) {
	s := &SlashCommand{Handler: &SingleCommandHandler{}}

	assert.Equal(t, 0, len(s.Handler.GetOptions()))

	o, err := s.AddOption("testOption", "Testing Option", true, discordgo.ApplicationCommandOptionBoolean)

	assert.Nil(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, &Option{&discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "testOption",
		Description: "Testing Option",
		Required:    true,
	}}, o)
	assert.Equal(t, 1, len(s.Handler.GetOptions()))
}

func TestAddChoice(t *testing.T) {
	o := &Option{&discordgo.ApplicationCommandOption{}}

	assert.Equal(t, 0, len(o.Choices))

	o.AddChoice("test", "5")

	assert.NotNil(t, o.Choices)
	assert.Equal(t, 1, len(o.Choices))
}
