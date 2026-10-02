package entity

import (
	"testing"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type doorCall struct {
	op  string
	m   *Map
	key DoorKey
}

type doorRecorder struct {
	MapSystem
	town  *Map
	calls []doorCall
}

func (r *doorRecorder) CreateReturnDoor(ch *Character, key DoorKey, skillID constant.SkillID) *Map {
	r.calls = append(r.calls, doorCall{op: "create", m: r.town, key: key})
	return r.town
}

func (r *doorRecorder) DespawnDoor(m *Map, key DoorKey, animated bool, notifyCounterpart bool) {
	r.calls = append(r.calls, doorCall{op: "despawn", m: m, key: key})
}

type doorWorld struct {
	GameWorld
	maps *doorRecorder
}

func (w doorWorld) GetMapSystem() MapSystem { return w.maps }

func newDoorTest() (*Character, *doorRecorder) {
	rec := &doorRecorder{town: &Map{id: 2}}
	ch := &Character{id: 7}
	ch.GameWorld = doorWorld{maps: rec}
	ch.Map = &Map{id: 1, Wz: &wz.Map{}}
	ch.Doors = NewDoorContainer(ch)
	return ch, rec
}

func mysticDoorBuff() *SkillBuff {
	return &SkillBuff{
		BaseBuff: &BaseBuff{},
		Wz:       &wz.Skill{ID: uint32(constant.SkillMysticDoor)},
	}
}

func (r *doorRecorder) expect(t *testing.T, want ...doorCall) {
	t.Helper()
	if len(r.calls) != len(want) {
		t.Fatalf("calls = %+v, want %+v", r.calls, want)
	}
	for i := range want {
		if r.calls[i] != want[i] {
			t.Fatalf("call %d = %+v, want %+v", i, r.calls[i], want[i])
		}
	}
}

func TestRecastBeforeResponseRemovesOldTownHalfFirst(t *testing.T) {
	ch, rec := newDoorTest()
	first := mysticDoorBuff()
	second := mysticDoorBuff()

	ch.Doors.Spawn(first)
	ch.Doors.Remove(first, true)
	ch.Doors.Spawn(second)

	rec.expect(t,
		doorCall{op: "create", m: rec.town, key: "7:1"},
		doorCall{op: "despawn", m: rec.town, key: "7:1"},
		doorCall{op: "create", m: rec.town, key: "7:2"},
	)
	if door := ch.Doors.SpawnField("7:1", constant.SkillMysticDoor, DoorEndpoint{}, DoorEndpoint{Map: ch.Map}); door != nil {
		t.Fatal("a late response of the first cast must not open a field door")
	}
}

func TestRecastWithoutUnbuffStillRemovesOldTownHalfFirst(t *testing.T) {
	ch, rec := newDoorTest()

	ch.Doors.Spawn(mysticDoorBuff())
	ch.Doors.Spawn(mysticDoorBuff())

	rec.expect(t,
		doorCall{op: "create", m: rec.town, key: "7:1"},
		doorCall{op: "despawn", m: rec.town, key: "7:1"},
		doorCall{op: "create", m: rec.town, key: "7:2"},
	)
}

func TestLateUnbuffKeepsNewerCast(t *testing.T) {
	ch, rec := newDoorTest()
	first := mysticDoorBuff()
	second := mysticDoorBuff()

	ch.Doors.Spawn(first)
	ch.Doors.Spawn(second)
	ch.Doors.Remove(first, true)

	rec.expect(t,
		doorCall{op: "create", m: rec.town, key: "7:1"},
		doorCall{op: "despawn", m: rec.town, key: "7:1"},
		doorCall{op: "create", m: rec.town, key: "7:2"},
	)
	if cast := ch.Doors.casts[constant.SkillMysticDoor]; cast == nil || cast.buff != second {
		t.Fatal("the newer cast must still wait for its town half")
	}
}
