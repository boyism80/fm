package entity

type CharacterListener interface {
	OnDialog(ch *Character, msg string, prev bool, next bool)
}
