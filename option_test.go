package harmonia

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestAddOption(t *testing.T) {

	s := NewSlashCommand("test")

	assert.Equal(t, 0, len(s.options))

	assert.Panics(t, func() { NewOption("", discordgo.ApplicationCommandOptionChannel) })
	opt := NewOption("testOption", discordgo.ApplicationCommandOptionBoolean).
		WithDescription("Testing Option").
		IsRequired()
	s.WithOptions(opt)

	assert.NotNil(t, opt)
	assert.Equal(t, &Option{&discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "testOption",
		Description: "Testing Option",
		Required:    true,
	}}, opt)
	assert.Equal(t, 1, len(s.options))
}

func TestAddChoice(t *testing.T) {
	o := NewOption("test", discordgo.ApplicationCommandOptionString)

	assert.Equal(t, 0, len(o.Choices))

	c := o.AddChoice("test", "5")

	assert.NotNil(t, o.Choices)
	assert.Equal(t, 1, len(o.Choices))
	assert.Equal(t, o.Choices[0], c)
}
