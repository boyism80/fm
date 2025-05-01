package data

type GameData struct {
	maps     map[uint32]*MapTemplate
	monsters map[uint32]*MonsterTemplate
	items    map[uint32]*ItemTemplate
}

// 생성자
func NewGameData() *GameData {
	return &GameData{
		maps:     map[uint32]*MapTemplate{},
		monsters: map[uint32]*MonsterTemplate{},
		items:    map[uint32]*ItemTemplate{},
	}
}

// Getter만 제공하여 불변성 유지
func (g *GameData) Map(id uint32) *MapTemplate {
	return g.maps[id]
}

func (g *GameData) Monster(id uint32) *MonsterTemplate {
	return g.monsters[id]
}

func (g *GameData) Item(id uint32) *ItemTemplate {
	return g.items[id]
}
