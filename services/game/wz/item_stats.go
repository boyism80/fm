package wz

type BasicStats struct {
	Str uint16
	Dex uint16
	Int uint16
	Luk uint16
}

type RequiredStats struct {
	BasicStats
	Level uint8
	Class int
}

type AbilityStats struct {
	BasicStats
	PAD       uint16
	MAD       uint16
	PDD       uint16
	MDD       uint16
	Speed     uint16
	Jump      uint16
	ACC       uint16
	EVA       int
	MaxHP     uint16
	MaxMP     uint16
	PVPDamage int

	Avoid uint16
	Hands uint16
}
