package constant

type CharacterRole uint8

const (
	RoleUser CharacterRole = iota
	RoleAdmin
)

func AllCharacterRoles() map[string]CharacterRole {
	return map[string]CharacterRole{
		"User":  RoleUser,
		"Admin": RoleAdmin,
	}
}
