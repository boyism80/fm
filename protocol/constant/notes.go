package constant

type NotesPacket uint8

const (
	NotesPacketList    NotesPacket = 3
	NotesPacketSent    NotesPacket = 4
	NotesPacketFailed  NotesPacket = 5
	NotesPacketRefresh NotesPacket = 7
)

type NoteKind uint8

const (
	NoteKindNormal     NoteKind = 0
	NoteKindPopulation NoteKind = 1
	NoteKindItem       NoteKind = 3
)
