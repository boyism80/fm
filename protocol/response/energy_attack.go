package response

type EnergyAttack struct {
	Attack
}

func (*EnergyAttack) Opcode() uint16 { return 0x86 }
