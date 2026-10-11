package response

import "github.com/boyism80/fm/stream"

type FamilyPrivileges struct {
}

func (p *FamilyPrivileges) Opcode() uint16 {
	return 0x4F
}

func (p *FamilyPrivileges) Serialize(w *stream.StreamWriter) error {
	names := []string{
		"패밀리원에게 바로 이동",
		"패밀리원 바로 소환",
		"나만의 드롭률 1.2배(15분)",
		"나만의 경험치 1.2배(15분)",
		"나만의 드롭률 1.2배(30분)",
		"나만의 경험치 1.2배(30분)",
		"패밀리원의 단결(30분)",
	}
	costs := []uint32{300, 500, 700, 800, 1000, 1200, 3000}
	kinds := []uint8{0, 1, 3, 2, 3, 2, 4}
	w.WriteU32(uint32(len(names)))
	for i, name := range names {
		w.WriteU8(kinds[i])
		w.WriteU32(costs[i])
		w.WriteU32(1)
		w.WriteStr16(name)
		w.WriteStr16(name + " (하루 1회)")
	}
	return nil
}

func (p *FamilyPrivileges) Deserialize(reader *stream.StreamReader) {
}
