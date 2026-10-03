package client

import (
	"testing"

	"github.com/boyism80/fm/services/game/entity"
)

func TestSetCharacterAfterLogoutFails(t *testing.T) {
	c := &GameClient{}
	if c.Logout() != nil {
		t.Fatal("no character yet")
	}
	if c.SetCharacter(&entity.Character{}) {
		t.Fatal("a login that finishes after logout must not set the character")
	}
	if c.GetCharacter() != nil {
		t.Fatal("character must stay unset")
	}
}

func TestLogoutTakesCharacter(t *testing.T) {
	c := &GameClient{}
	ch := &entity.Character{}
	if c.SetCharacter(ch) == false {
		t.Fatal("SetCharacter before logout must succeed")
	}
	if c.Logout() != ch {
		t.Fatal("logout must return the set character")
	}
	if c.GetCharacter() != nil {
		t.Fatal("logout must take the character off the client")
	}
}
