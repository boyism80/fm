package response

const (
	buddyNameFieldLen  = 13
	buddyGroupFieldLen = 17
)

type BuddyEntry struct {
	CharacterID uint32
	Name        string
	Pending     bool
	Channel     int32
	Group       string
}
