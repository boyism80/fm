package constant

type MistType uint8

const (
	MistTypePoison MistType = 0
	MistTypeSmoke  MistType = 2
)

func AllMistTypes() map[string]MistType {
	return map[string]MistType{
		"Poison": MistTypePoison,
		"Smoke":  MistTypeSmoke,
	}
}
