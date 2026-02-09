// Package constant defines buff types for temporary buffs.
// Mask is the client bit flag; Position is the mask group index (1-based, 2–4).
// MaxBuffFlag = 4: packet uses 4 int32 mask words.
package constant

const MaxBuffFlag = 4

// BuffFlag identifies a buff type for packets (mask value + position).
type BuffFlag struct {
	Mask     uint32
	Position int
}

// Common buff flags (position 4 = first mask group).
var (
	BuffFlagWATK          = BuffFlag{0x1, 4}
	BuffFlagWDEF          = BuffFlag{0x2, 4}
	BuffFlagMATK          = BuffFlag{0x4, 4}
	BuffFlagMDEF          = BuffFlag{0x8, 4}
	BuffFlagACC           = BuffFlag{0x10, 4}
	BuffFlagAVOID         = BuffFlag{0x20, 4}
	BuffFlagHANDS         = BuffFlag{0x40, 4}
	BuffFlagSPEED         = BuffFlag{0x80, 4}
	BuffFlagJUMP          = BuffFlag{0x100, 4}
	BuffFlagMAGIC_GUARD   = BuffFlag{0x200, 4}
	BuffFlagDARKSIGHT     = BuffFlag{0x400, 4}
	BuffFlagBOOSTER       = BuffFlag{0x800, 4}
	BuffFlagPOWERGUARD    = BuffFlag{0x1000, 4}
	BuffFlagMAXHP         = BuffFlag{0x2000, 4}
	BuffFlagMAXMP         = BuffFlag{0x4000, 4}
	BuffFlagINVINCIBLE    = BuffFlag{0x8000, 4}
	BuffFlagSOULARROW     = BuffFlag{0x10000, 4}
	BuffFlagCOMBO         = BuffFlag{0x200000, 4}
	BuffFlagWK_CHARGE     = BuffFlag{0x400000, 4}
	BuffFlagDRAGONBLOOD   = BuffFlag{0x800000, 4}
	BuffFlagHOLY_SYMBOL   = BuffFlag{0x1000000, 4}
	BuffFlagMESOUP        = BuffFlag{0x2000000, 4}
	BuffFlagSHADOWPARTNER = BuffFlag{0x4000000, 4}
	BuffFlagPICKPOCKET    = BuffFlag{0x8000000, 4}
	BuffFlagMESOGUARD     = BuffFlag{0x10000000, 4}
	BuffFlagHP_LOSS_GUARD = BuffFlag{0x20000000, 4}
)

// Position 3.
var (
	BuffFlagMORPH          = BuffFlag{0x2, 3}
	BuffFlagRECOVERY       = BuffFlag{0x4, 3}
	BuffFlagMapleWarrior   = BuffFlag{0x8, 3}
	BuffFlagSTANCE         = BuffFlag{0x10, 3}
	BuffFlagSHARP_EYES     = BuffFlag{0x20, 3}
	BuffFlagManaReflection = BuffFlag{0x40, 3}
	BuffFlagSPIRIT_CLAW    = BuffFlag{0x100, 3}
	BuffFlagINFINITY       = BuffFlag{0x200, 3}
	BuffFlagBERSERK_FURY   = BuffFlag{0x2000000, 3}
	BuffFlagDIVINE_BODY    = BuffFlag{0x4000000, 3}
	BuffFlagFINALATTACK    = BuffFlag{0x20000000, 3}
)

// Position 2.
var (
	BuffFlagENERGY_CHARGE  = BuffFlag{0x2, 2}
	BuffFlagDASH_SPEED     = BuffFlag{0x4, 2}
	BuffFlagDASH_JUMP      = BuffFlag{0x8, 2}
	BuffFlagMONSTER_RIDING = BuffFlag{0x10, 2}
	BuffFlagSPEED_INFUSION = BuffFlag{0x20, 2}
	BuffFlagHOMING_BEACON  = BuffFlag{0x40, 2}
	BuffFlagEXPRATE        = BuffFlag{0x400, 2}
	BuffFlagDROP_RATE      = BuffFlag{0x800, 2}
	BuffFlagMESO_RATE      = BuffFlag{0x1000, 2}
)
