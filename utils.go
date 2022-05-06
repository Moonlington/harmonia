package harmonia

import "github.com/bwmarrin/discordgo"

// ParseComponentMatrix parses a 2D slice of MessageComponents and returns a 1D slice of MessageComponents with ActionsRows.
func ParseComponentMatrix(components [][]discordgo.MessageComponent) []discordgo.MessageComponent {
	comp := make([]discordgo.MessageComponent, len(components))
	for i, c := range components {
		comp[i] = &discordgo.ActionsRow{Components: c}
	}
	return comp
}
