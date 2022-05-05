package harmonia

import "github.com/bwmarrin/discordgo"

func ParseComponentMatrix(components [][]discordgo.MessageComponent) []discordgo.MessageComponent {
	comp := make([]discordgo.MessageComponent, len(components))
	for i, c := range components {
		comp[i] = &discordgo.ActionsRow{Components: c}
	}
	return comp
}
