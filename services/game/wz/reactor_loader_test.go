package wz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boyism80/fm/services/game/constant"
)

func TestLoadReactor(t *testing.T) {
	path := filepath.Join("..", "..", "resources", "wz", "Reactor.wz", "02111001.img.xml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Test file not found: %s", path)
		return
	}

	reactor, err := loadReactor(path)
	if err != nil {
		t.Fatalf("Failed to load reactor: %v", err)
	}
	if reactor == nil {
		t.Fatal("loadReactor returned nil")
	}
	if reactor.ID != 2111001 {
		t.Errorf("expected id 2111001, got %d", reactor.ID)
	}
	if len(reactor.States) == 0 {
		t.Fatal("expected at least one state")
	}

	state0 := reactor.States[0]
	if state0 == nil {
		t.Fatal("expected state 0 event")
	}
	if state0.Type != constant.ReactorEventTypeItem {
		t.Errorf("expected state 0 type Item, got %s", state0.Type)
	}
}

func TestLoadReactorRawWZNoSyntheticEvent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "0100000.img.xml")
	xmlData := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<imgdir name="0100000.img">
    <imgdir name="info">
        <int name="activateByTouch" value="0"/>
    </imgdir>
    <imgdir name="0">
        <imgdir name="event">
            <imgdir name="0">
                <int name="type" value="0"/>
                <int name="state" value="1"/>
            </imgdir>
        </imgdir>
    </imgdir>
    <imgdir name="1">
    </imgdir>
</imgdir>`
	if err := os.WriteFile(path, []byte(xmlData), 0o644); err != nil {
		t.Fatalf("write test xml: %v", err)
	}

	reactor, err := loadReactor(path)
	if err != nil {
		t.Fatalf("Failed to load reactor: %v", err)
	}

	state1 := reactor.States[1]
	if state1 != nil {
		t.Fatal("state without event/0 must not get a synthetic event")
	}

	state0 := reactor.States[0]
	if state0 == nil {
		t.Fatal("expected state 0 event")
	}
	if state0.Type != constant.ReactorEventTypeHit {
		t.Errorf("expected type Hit, got %s", state0.Type)
	}
	if state0.NextState != 1 {
		t.Errorf("expected next state 1, got %d", state0.NextState)
	}
}

func TestLoadReactorTerminalTypeNormalizedToNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "0100000.img.xml")
	xmlData := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<imgdir name="0100000.img">
    <imgdir name="0">
        <imgdir name="event">
            <imgdir name="0">
                <int name="type" value="0"/>
                <int name="state" value="1"/>
            </imgdir>
        </imgdir>
    </imgdir>
    <imgdir name="1">
        <imgdir name="event">
            <imgdir name="0">
                <int name="type" value="999"/>
                <int name="state" value="1"/>
            </imgdir>
        </imgdir>
    </imgdir>
</imgdir>`
	if err := os.WriteFile(path, []byte(xmlData), 0o644); err != nil {
		t.Fatalf("write test xml: %v", err)
	}

	reactor, err := loadReactor(path)
	if err != nil {
		t.Fatalf("Failed to load reactor: %v", err)
	}

	if reactor.States[1] != nil {
		t.Fatal("WZ terminal type 999 must normalize to nil event")
	}
}

func TestResourcesGetReactorLink(t *testing.T) {
	linked := &Reactor{
		ID: 9908001,
		States: map[byte]*ReactorEvent{
			0: {Type: constant.ReactorEventTypeItem},
		},
	}
	alias := &Reactor{
		ID:   9908002,
		Info: ReactorInfo{Link: 9908001},
	}
	resources := &Resources{
		Reactors: map[uint32]*Reactor{
			9908001: linked,
			9908002: alias,
		},
	}

	got := resources.GetReactor(9908002)
	if got != linked {
		t.Fatal("GetReactor should resolve info.link to linked reactor def")
	}
}
