package dashboard

import (
	botactor "github.com/boyism80/fm/bot/actor"
)

type CatalogSuite struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	Bots      int    `json:"bots"`
	Serial    bool   `json:"serial"`
	Explicit  bool   `json:"explicit"`
	TimeoutMs int    `json:"timeout_ms"`
	Error     string `json:"error,omitempty"`
}

func (s *Server) catalog() ([]*CatalogSuite, error) {
	suites, err := botactor.LoadSuites(s.cfg.ScriptDir)
	if err != nil {
		return nil, err
	}

	result := make([]*CatalogSuite, 0, len(suites))
	for _, suite := range suites {
		entry := &CatalogSuite{Name: suite.Name, Serial: suite.Serial, Explicit: suite.Explicit, TimeoutMs: s.cfg.SuiteTimeoutMs}
		if suite.TimeoutMs > 0 {
			entry.TimeoutMs = suite.TimeoutMs
		}
		entry.Title, entry.Bots, err = suite.Describe(s.cfg, s.wz)
		if err != nil {
			entry.Error = err.Error()
		}
		result = append(result, entry)
	}
	return result, nil
}
