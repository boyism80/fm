package response

import "github.com/boyism80/fm/stream"

type PartyMemberStatus struct {
	CharacterID uint32
	Name        string
	Class       uint32
	Level       uint32
	Channel     int32
	MapID       uint32
	DoorTown    uint32
	DoorTarget  uint32
	DoorX       int32
	DoorY       int32
}

func normalizePartyMembersSix(members []PartyMemberStatus) []PartyMemberStatus {
	out := make([]PartyMemberStatus, 6)
	for i := 0; i < len(out) && i < len(members); i++ {
		out[i] = members[i]
	}
	return out
}

func writePartyStatusBlock(w *stream.StreamWriter, forChannel int32, leaderCharacterID uint32, members []PartyMemberStatus, leaving bool) {
	slots := normalizePartyMembersSix(members)
	for _, m := range slots {
		w.WriteU32(m.CharacterID)
	}
	for _, m := range slots {
		w.WriteStaticStr(m.Name, 13)
	}
	for _, m := range slots {
		w.WriteU32(m.Class)
	}
	for _, m := range slots {
		w.WriteU32(m.Level)
	}
	for _, m := range slots {
		w.Write32(m.Channel)
	}
	w.WriteU32(leaderCharacterID)
	for _, m := range slots {
		if m.Channel == forChannel {
			w.WriteU32(m.MapID)
		} else {
			w.WriteU32(0)
		}
	}
	for _, m := range slots {
		if m.Channel == forChannel && !leaving {
			w.WriteU32(m.DoorTown)
			w.WriteU32(m.DoorTarget)
			w.Write32(m.DoorX)
			w.Write32(m.DoorY)
			continue
		}
		if leaving {
			w.WriteU32(999999999)
			w.WriteU32(999999999)
			w.Write64(-1)
			continue
		}
		w.WriteU32(0)
		w.WriteU32(0)
		w.Write64(0)
	}
}
