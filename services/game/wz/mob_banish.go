package wz

type MobBanish struct {
	Message string
	MapID   int32
	Portal  string
}

func parseMobBanish(banNode *node) *MobBanish {
	if banNode == nil {
		return nil
	}
	mapID := -1
	portal := "sp"
	for _, ch := range banNode.Children {
		if ch.Name == "banMap" {
			for _, entry := range ch.Children {
				if entry.Name != "0" {
					continue
				}
				mapID = nodeInt(&entry, "field", -1)
				for _, sf := range entry.Strings {
					if sf.Name == "portal" {
						portal = sf.Value
					}
				}
			}
		}
	}
	msg := ""
	for _, sf := range banNode.Strings {
		if sf.Name == "banMsg" {
			msg = sf.Value
			break
		}
	}
	if mapID < 0 {
		return nil
	}
	return &MobBanish{
		Message: msg,
		MapID:   int32(mapID),
		Portal:  portal,
	}
}
