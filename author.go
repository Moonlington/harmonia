package harmonia

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// An Author describes either a User or Member, depending if the message was sent in a Guild or DMs.
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

// AuthorFromUser returns an Author from a *discordgo.User.
func AuthorFromUser(user *discordgo.User) *Author {
	return &Author{User: user, IsMember: false, Nick: user.GlobalName}
}

// AuthorFromInteraction uses the information obtained from the Interaction to create an Author.
func AuthorFromInteraction(h *Harmonia, i *discordgo.Interaction) (a *Author, err error) {
	if i.Member == nil {
		return AuthorFromUser(i.User), nil
	}

	i.Member.GuildID = i.GuildID
	return AuthorFromMember(h, i.Member)
}

// AuthorFromMember returns an Author from a *discordgo.Member.
func AuthorFromMember(h *Harmonia, member *discordgo.Member) (*Author, error) {
	guild, err := h.Guild(member.GuildID)
	if err != nil {
		return nil, err
	}

	roles, err := RolesFromMember(h, member)
	if err != nil {
		return nil, err
	}

	a := &Author{User: member.User,
		IsMember:     true,
		Guild:        guild,
		JoinedAt:     member.JoinedAt,
		Nick:         member.Nick,
		Deaf:         member.Deaf,
		Mute:         member.Mute,
		Roles:        roles,
		PremiumSince: member.PremiumSince,
	}
	a.Avatar = member.Avatar
	return a, nil
}

// RolesFromMember returns a slice of *discordgo.Role from a *discordgo.Member.
func RolesFromMember(h *Harmonia, member *discordgo.Member) ([]*discordgo.Role, error) {
	guildroles, err := h.GuildRoles(member.GuildID)
	if err != nil {
		return nil, err
	}

	roles := make([]*discordgo.Role, 0, len(member.Roles))
	for _, roleid := range member.Roles {
		for _, role := range guildroles {
			if role.ID == roleid {
				roles = append(roles, role)
			}
		}
	}

	return roles, nil
}
