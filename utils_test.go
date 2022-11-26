package harmonia

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestParseComponentMatrix(t *testing.T) {
	components := [][]discordgo.MessageComponent{
		{
			discordgo.Button{Label: "Button 1", CustomID: "button1"},
			discordgo.Button{Label: "Button 2", CustomID: "button2"},
		}, {
			discordgo.Button{Label: "Button 3", CustomID: "button3"},
		},
	}
	parsedMatrix := ParseComponentMatrix(components)
	correctMatrix := []discordgo.MessageComponent{&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{Label: "Button 1", CustomID: "button1"},
		discordgo.Button{Label: "Button 2", CustomID: "button2"},
	}}, &discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{Label: "Button 3", CustomID: "button3"},
	}}}

	assert.Equal(t, correctMatrix, parsedMatrix)
}
