package entity

type SavedLocations struct {
	owner   *Character
	entries map[string]uint32
}

func (s *SavedLocations) Save(name string) {
	if name == "" {
		return
	}
	m := s.owner.GetMap()
	if m == nil {
		return
	}
	s.entries[name] = m.TemplateID()
}

func (s *SavedLocations) Find(name string) (uint32, bool) {
	if name == "" {
		return 0, false
	}
	mapID, ok := s.entries[name]
	return mapID, ok
}

func (s *SavedLocations) Clear(name string) {
	if name == "" {
		return
	}
	delete(s.entries, name)
}
