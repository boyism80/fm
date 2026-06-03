package dto

type AllianceInfo struct {
	AllianceID uint32
	Name       string
	RankTitles [5]string
	GuildIDs   []uint32
	Capacity   uint32
	Notice     string
}

type AllianceMembershipChangeMember struct {
	CharacterID  uint32
	AllianceRank uint8
}

type AllianceMembershipChangeGuild struct {
	GuildID uint32
	Members []AllianceMembershipChangeMember
}

type AllianceGuildMemberRank struct {
	CharacterID  uint32
	AllianceRank uint8
}
