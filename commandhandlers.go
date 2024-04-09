package harmonia

// // A CommandGroupHandler describes the handler for a command with subcommands.
// type CommandGroupHandler struct {
// 	// Subcommands contains a map of SubSlashCommands
// 	Subcommands map[string]*SlashSubcommand
// }

// // Do handles the subcommands given to the SubcommandHandler
// func (s *CommandGroupHandler) Do(h *Harmonia, i *Invocation) {
// 	options := i.options
// 	if sc, ok := s.Subcommands[options[0].Name]; ok {
// 		i.options = options[0].Options
// 		sc.Handler.Do(h, i)
// 	}
// }

// // GetOptions returns the subcommands parsed as ApplicationCommandOptions
// func (s *CommandGroupHandler) GetOptions() []*discordgo.ApplicationCommandOption {
// 	options := make([]*discordgo.ApplicationCommandOption, len(s.Subcommands))
// 	i := 0
// 	for _, sc := range s.Subcommands {
// 		t := discordgo.ApplicationCommandOptionSubCommand
// 		if sc.IsGroup {
// 			t = discordgo.ApplicationCommandOptionSubCommandGroup
// 		}

// 		options[i] = &discordgo.ApplicationCommandOption{
// 			Name:        sc.Name,
// 			Description: sc.Description,
// 			Options:     sc.Handler.GetOptions(),
// 			Type:        t,
// 		}
// 		i++
// 	}
// 	return options
// }
