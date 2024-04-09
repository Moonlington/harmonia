package harmonia

import (
	"log"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	h *Harmonia

	envBotToken string
	// envGuild    = os.Getenv("H_GUILD")
	// envAdmin    = os.Getenv("H_ADMIN")
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	envBotToken = os.Getenv("H_TOKEN")

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

func TestAddCommand(t *testing.T) {
	harm := &Harmonia{Commands: make(map[string]CommandHandler)}

	command := NewSlashCommand("test")
	t.Run("Correct Slash Command", func(t *testing.T) {
		err := harm.AddCommand(command)
		assert.Nil(t, err)
		assert.Equal(t, command, harm.Commands["test"])
	})
	t.Run("Duplicate Slash Command", func(t *testing.T) {
		err := harm.AddCommand(command)
		assert.EqualError(t, err, "command 'test' already exists")
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
