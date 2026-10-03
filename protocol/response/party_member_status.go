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
	n := len(members)
	if n > len(out) {
		n = len(out)
	}
	for i := 0; i < n; i++ {
		out[i] = members[i]
	}
	for i := n; i < len(out); i++ {
		out[i].Channel = -2
	}
	return out
}

func readPartyStatusBlock(r *stream.StreamReader) (uint32, []PartyMemberStatus) {
	slots := make([]PartyMemberStatus, 6)
	for i := range slots {
		slots[i].CharacterID = r.ReadU32()
	}
	for i := range slots {
		slots[i].Name = r.ReadStaticStr(13)
	}
	for i := range slots {
		slots[i].Class = r.ReadU32()
	}
	for i := range slots {
		slots[i].Level = r.ReadU32()
	}
	for i := range slots {
		slots[i].Channel = r.Read32()
	}
	leaderCharacterID := r.ReadU32()
	for i := range slots {
		slots[i].MapID = r.ReadU32()
	}
	for i := range slots {
		slots[i].DoorTown = r.ReadU32()
		slots[i].DoorTarget = r.ReadU32()
		slots[i].DoorX = r.Read32()
		slots[i].DoorY = r.Read32()
	}

	var members []PartyMemberStatus
	for _, m := range slots {
		if m.CharacterID != 0 {
			members = append(members, m)
		}
	}
	return leaderCharacterID, members
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
		} else {
			w.WriteU32(0)
			w.WriteU32(0)
			w.Write64(0)
		}
	}
}
