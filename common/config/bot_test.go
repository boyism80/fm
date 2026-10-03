package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBotDefaultsAndValidate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bot.yaml")
	if err := os.WriteFile(path, []byte("channel: 0\nworld_id: 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := LoadBot(path)
	if err != nil {
		t.Fatal(err)
	}
	if b.Login.Host != "127.0.0.1" || b.Login.Port != 8484 {
		t.Fatalf("login=%s:%d", b.Login.Host, b.Login.Port)
	}
	if b.Channel != 0 || b.WorldID != 0 || b.Seats != 3 {
		t.Fatalf("channel=%d world=%d seats=%d", b.Channel, b.WorldID, b.Seats)
	}
	if b.Password != "bot-pass" || b.TimeoutMs != 10000 || b.SuiteTimeoutMs != 180000 {
		t.Fatalf("password=%q timeout=%d suite=%d", b.Password, b.TimeoutMs, b.SuiteTimeoutMs)
	}
	if b.ScriptDir != "script/integration" || b.WzPath != "resources/wz" {
		t.Fatalf("script=%q wz=%q", b.ScriptDir, b.WzPath)
	}
}

func TestLoadBotRejectsInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bot.yaml")
	if err := os.WriteFile(path, []byte("channel: -1\nseats: 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadBot(path); err == nil {
		t.Fatal("expected error")
	}
}
