package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/boyism80/fm/services/game/entity"
)

type StateMachineRegistry struct {
	gs     *GameServer
	mu     sync.RWMutex
	groups map[string]*entity.StateMachineGroup
}

func NewStateMachineRegistry(gs *GameServer) *StateMachineRegistry {
	return &StateMachineRegistry{
		gs:     gs,
		groups: make(map[string]*entity.StateMachineGroup),
	}
}

func (r *StateMachineRegistry) Get(name string) *entity.StateMachineGroup {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.groups[name]
}

func (r *StateMachineRegistry) Register(group *entity.StateMachineGroup) {
	if r == nil || group == nil || group.Name == "" {
		return
	}
	r.mu.Lock()
	r.groups[group.Name] = group
	r.mu.Unlock()
}

func (r *StateMachineRegistry) LoadFromScripts(dir string) error {
	if r == nil || r.gs == nil {
		return fmt.Errorf("state machine registry not ready")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("state machine: script dir %s not found, skipping", dir)
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".lua") {
			continue
		}
		stem := strings.TrimSuffix(name, ".lua")
		scriptPath := filepath.ToSlash(filepath.Join(dir, name))
		group := entity.NewStateMachineGroup(stem, scriptPath, r.gs)
		r.Register(group)
		group.Init()
		log.Printf("state machine: loaded group %s (%s)", stem, scriptPath)
	}
	return nil
}

func (r *StateMachineRegistry) CancelAll() {
	if r == nil {
		return
	}
	r.mu.RLock()
	groups := make([]*entity.StateMachineGroup, 0, len(r.groups))
	for _, g := range r.groups {
		groups = append(groups, g)
	}
	r.mu.RUnlock()
	for _, g := range groups {
		g.CancelSchedule()
	}
}
