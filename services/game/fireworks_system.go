package game

import "sync"

// FireWorksSystem keeps the channel's shared donation progress (old/src/server/maps/AramiaFireWorks.java).
// It lives in memory only, so a restart begins again from a sixth of each goal.
type FireWorksSystem struct {
	mu       sync.Mutex
	goals    map[string]int
	progress map[string]int
}

func NewFireWorksSystem() *FireWorksSystem {
	goals := map[string]int{"keg": 2400, "sun": 3000}
	progress := make(map[string]int, len(goals))
	for kind, goal := range goals {
		progress[kind] = goal / 6
	}
	return &FireWorksSystem{goals: goals, progress: progress}
}

func (s *FireWorksSystem) Start() error { return nil }
func (s *FireWorksSystem) Stop()        {}

// Percentage is the progress in hundredths of a percent (10000 is full), or -1 for an unknown kind.
func (s *FireWorksSystem) Percentage(kind string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	goal, ok := s.goals[kind]
	if !ok {
		return -1
	}
	return s.progress[kind] * 10000 / goal
}

// Give adds a donation and reports whether it filled the goal; a full goal starts over from zero.
func (s *FireWorksSystem) Give(kind string, count int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	goal, ok := s.goals[kind]
	if !ok || count <= 0 {
		return false
	}
	s.progress[kind] += count
	if s.progress[kind] < goal {
		return false
	}
	s.progress[kind] = 0
	return true
}
