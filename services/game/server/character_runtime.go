package server

import (
	"fmt"
	"strings"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

type ServerCharacterRuntime struct {
	gs     *GameServer
	mu     sync.RWMutex
	byCID  map[uint32]*characterRuntimeEntry
	byName map[string]uint32
}

type characterRuntimeEntry struct {
	name      string
	mapPID    *actor.PID
	ensureBuf []*g_actor.EnsureDeliver
}

func NewServerCharacterRuntime(gs *GameServer) *ServerCharacterRuntime {
	return &ServerCharacterRuntime{
		gs:     gs,
		byCID:  make(map[uint32]*characterRuntimeEntry),
		byName: make(map[string]uint32),
	}
}

func normalizeRuntimeName(name string) string {
	return strings.TrimSpace(name)
}

func (r *ServerCharacterRuntime) RegisterCharacter(characterID uint32, name string) error {
	if r == nil {
		return fmt.Errorf("runtime: nil receiver")
	}
	if characterID == 0 {
		return fmt.Errorf("runtime: invalid character id 0")
	}
	key := normalizeRuntimeName(name)
	if key == "" {
		return fmt.Errorf("runtime: empty character name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byCID[characterID]; exists {
		return fmt.Errorf("runtime: character %d already registered", characterID)
	}
	if existingID, taken := r.byName[key]; taken && existingID != characterID {
		return fmt.Errorf("runtime: name %q already in use", key)
	}
	r.byCID[characterID] = &characterRuntimeEntry{name: key}
	r.byName[key] = characterID
	return nil
}

func (r *ServerCharacterRuntime) UnregisterCharacter(characterID uint32) {
	if r == nil || characterID == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ent, ok := r.byCID[characterID]
	if !ok {
		return
	}
	delete(r.byName, ent.name)
	delete(r.byCID, characterID)
}

func (r *ServerCharacterRuntime) SetMapPID(characterID uint32, pid *actor.PID) error {
	if r == nil {
		return fmt.Errorf("runtime: nil receiver")
	}
	if characterID == 0 {
		return fmt.Errorf("runtime: invalid character id 0")
	}
	r.mu.Lock()
	ent, ok := r.byCID[characterID]
	if !ok {
		r.mu.Unlock()
		return fmt.Errorf("runtime: character %d not registered", characterID)
	}
	ent.mapPID = pid
	r.mu.Unlock()

	if pid != nil && r.gs != nil {
		r.gs.flushEnsureBuffer(characterID, pid)
	}
	return nil
}

func (r *ServerCharacterRuntime) Exists(characterID uint32) bool {
	if r == nil || characterID == 0 {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byCID[characterID]
	return ok
}

func (r *ServerCharacterRuntime) GetMapPID(characterID uint32) (pid *actor.PID, ok bool) {
	if r == nil || characterID == 0 {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	ent, exists := r.byCID[characterID]
	if !exists {
		return nil, false
	}
	return ent.mapPID, true
}

func (r *ServerCharacterRuntime) GetCharacterIDByName(name string) (uint32, bool) {
	if r == nil {
		return 0, false
	}
	key := normalizeRuntimeName(name)
	if key == "" {
		return 0, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byName[key]
	return id, ok
}

func (r *ServerCharacterRuntime) enqueueEnsure(characterID uint32, d *g_actor.EnsureDeliver) error {
	if r == nil || d == nil {
		return fmt.Errorf("runtime: nil receiver or deliver")
	}
	if characterID == 0 {
		return fmt.Errorf("runtime: invalid character id 0")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ent, ok := r.byCID[characterID]
	if !ok {
		return fmt.Errorf("runtime: character %d not registered", characterID)
	}
	ent.ensureBuf = append(ent.ensureBuf, d)
	return nil
}

func (r *ServerCharacterRuntime) takeEnsureBuffer(characterID uint32) []*g_actor.EnsureDeliver {
	if r == nil || characterID == 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ent, ok := r.byCID[characterID]
	if !ok || len(ent.ensureBuf) == 0 {
		return nil
	}
	buf := ent.ensureBuf
	ent.ensureBuf = nil
	return buf
}

func (r *ServerCharacterRuntime) removeEnsureFromBuffer(characterID uint32, correlationID uint64) {
	if r == nil || characterID == 0 || correlationID == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	ent, ok := r.byCID[characterID]
	if !ok || len(ent.ensureBuf) == 0 {
		return
	}
	out := ent.ensureBuf[:0]
	for _, d := range ent.ensureBuf {
		if d == nil || d.CorrelationID != correlationID {
			out = append(out, d)
		}
	}
	ent.ensureBuf = out
}
