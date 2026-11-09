package wz

// partyPlayMobs are the Ellin Forest party hunting mobs the original server read with its own HP and EXP, as bosses
// that still show their HP and earn the map's party bonus.
var partyPlayMobs = map[uint32]struct {
	maxHP int
	exp   uint32
}{
	4250000: {16000, 470},
	4250001: {20000, 560},
	5250000: {21000, 580},
	5250001: {22000, 640},
	5250002: {23000, 700},
}

// partyPlayMaps take a 25% party bonus; any other map takes half its WZ partyBonusR.
var partyPlayMaps = map[uint32]bool{
	300010000: true,
	300010100: true,
	300010200: true,
	300010300: true,
	300010400: true,
	300020000: true,
	300020100: true,
	300020200: true,
	300030000: true,
}

const partyPlayRate = 25
