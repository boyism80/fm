package entity

func (ch *Character) SaveLocation(name string) {
	if name == "" {
		return
	}
	m := ch.GetMap()
	if m == nil {
		return
	}
	ch.savedLocations[name] = m.GetMapID()
}

func (ch *Character) SavedLocation(name string) (uint32, bool) {
	if name == "" {
		return 0, false
	}
	mapID, ok := ch.savedLocations[name]
	return mapID, ok
}

func (ch *Character) ClearSavedLocation(name string) {
	if name == "" {
		return
	}
	delete(ch.savedLocations, name)
}
