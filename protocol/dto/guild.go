package dto

type GuildMemberStatus struct {
	CharacterID  uint32
	Name         string
	ClassID      uint32
	Level        uint32
	GuildRank    uint32
	Online       bool
	AllianceRank *uint32
}

type GuildInfo struct {
	GuildID     uint32
	Name        string
	RankTitles  [5]string
	Members     []GuildMemberStatus
	Capacity    uint32
	LogoBG      uint16
	LogoBGColor uint8
	Logo        uint16
	LogoColor   uint8
	Notice      string
	GP          uint32
	AllianceID  *uint32
}

type GuildRankingEntry struct {
	Name        string
	GP          uint32
	Logo        uint32
	LogoColor   uint32
	LogoBG      uint32
	LogoBGColor uint32
}
