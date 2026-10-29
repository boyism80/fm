package cluster

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

// catalog reads the scripts on the master's disk; each agent runs a suite from its own copy.
func (m *Master) catalog() ([]*CatalogSuite, error) {
	suites, err := botactor.LoadSuites(m.cfg.ScriptDir)
	if err != nil {
		return nil, err
	}

	result := make([]*CatalogSuite, 0, len(suites))
	for _, suite := range suites {
		entry := &CatalogSuite{Name: suite.Name, Serial: suite.Serial, Explicit: suite.Explicit, TimeoutMs: m.cfg.SuiteTimeoutMs}
		if suite.TimeoutMs > 0 {
			entry.TimeoutMs = suite.TimeoutMs
		}
		entry.Title, entry.Bots, err = suite.Describe(m.cfg, m.wz)
		if err != nil {
			entry.Error = err.Error()
		}
		result = append(result, entry)
	}
	return result, nil
}
