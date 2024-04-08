package harmonia

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

var (
	h *Harmonia

	envBotToken = os.Getenv("H_TOKEN")
	// envGuild    = os.Getenv("H_GUILD")
	// envAdmin    = os.Getenv("H_ADMIN")
)

func TestMain(m *testing.M) {
	if envBotToken != "" {
		if harm, err := New(envBotToken); err == nil {
			h = harm
		}
	}

	os.Exit(m.Run())
}

// TestNewToken tests the New() function with a Token.
func TestNewToken(t *testing.T) {
	if envBotToken == "" {
		t.Skip("Skipping New(token), H_TOKEN not set")
	}

	harm, err := New(envBotToken)

	assert.Nil(t, err)
	assert.NotNil(t, harm)
	assert.NotEqual(t, "", harm.Token)
}

func TestAddSlashCommand(t *testing.T) {
	harm := &Harmonia{Commands: make(map[string]*SlashCommand)}
	t.Run("Slash Command with no name", func(t *testing.T) {
		s, err := harm.AddSlashCommand("", "", func(h *Harmonia, i *Invocation) {})
		assert.Nil(t, s)
		assert.EqualError(t, err, "Empty Slash Command name")
	})
	t.Run("Slash Command with Invalid Name", func(t *testing.T) {
		s, err := harm.AddSlashCommand("test/", "", func(h *Harmonia, i *Invocation) {})
		assert.Nil(t, s)
		assert.EqualError(t, err, "Slash Command name does not match with the CHAT_INPUT regex.")
	})
	t.Run("Correct Slash Command", func(t *testing.T) {
		handlerFunc := func(h *Harmonia, i *Invocation) {}
		s, err := harm.AddSlashCommand("test", "", handlerFunc)
		assert.Nil(t, err)
		assert.NotNil(t, s)
		assert.NotNil(t, harm.Commands["test"])
	})
	t.Run("Duplicate Slash Command", func(t *testing.T) {
		s, err := harm.AddSlashCommand("test", "", func(h *Harmonia, i *Invocation) {})
		assert.Nil(t, s)
		assert.EqualError(t, err, "Slash Command 'test' already exists")
	})
}

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
