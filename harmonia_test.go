package harmonia

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	h *Harmonia

	envBotToken = os.Getenv("H_TOKEN")
	envGuild    = os.Getenv("H_GUILD")
	envAdmin    = os.Getenv("H_ADMIN")
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