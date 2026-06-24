package entity

func (m *Map) GetProperty(key string) (interface{}, bool) {
	if m == nil || key == "" {
		return nil, false
	}
	m.propertyMutex.RLock()
	defer m.propertyMutex.RUnlock()
	if m.properties == nil {
		return nil, false
	}
	val, ok := m.properties[key]
	return val, ok
}

func (m *Map) SetProperty(key string, value interface{}) {
	if m == nil || key == "" {
		return
	}
	m.propertyMutex.Lock()
	defer m.propertyMutex.Unlock()
	if m.properties == nil {
		m.properties = make(map[string]interface{})
	}
	m.properties[key] = value
}

func (m *Map) ClearProperties() {
	if m == nil {
		return
	}
	m.propertyMutex.Lock()
	defer m.propertyMutex.Unlock()
	m.properties = nil
}
