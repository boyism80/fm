package constant

// CharacterRole represents the character's permission role.
// Values are ordered: higher value means higher privilege. Use HasRoleAtLeast to compare.
type CharacterRole uint8

const (
	RoleUser  CharacterRole = iota // 일반유저
	RoleAdmin                      // 관리자
)
