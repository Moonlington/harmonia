package harmonia

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Name         string
	Description  string
	Handler      func(h *Harmonia, i *Invocation)
	Options      []*discordgo.ApplicationCommandOption
	GuildID      string
	registration *discordgo.ApplicationCommand
	// Subcommands map[string]*SlashCommand
}

type Invocation struct {
	*discordgo.Interaction
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *Author
}

type Author struct {
	*discordgo.User
	IsMember     bool
	Guild        *discordgo.Guild
	JoinedAt     time.Time
	Nick         string
	Deaf         bool
	Mute         bool
	Roles        []*discordgo.Role
	PremiumSince *time.Time
}
