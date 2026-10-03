package actor

import (
	"fmt"
	"time"
)

// TODO(M4): replace with Lua suites from script/integration/registry.lua.
const henesys = 100000000

var Suites = []*Suite{
	{
		Name:     "login_smoke",
		BotCount: 1,
		Scenarios: []Scenario{
			{
				Name: "instance_move",
				Run: func(a *SuiteActor, done func(bool)) {
					a.InstanceMove(a.bots[0], henesys, done)
				},
			},
		},
	},
	{
		Name:     "seat_pair",
		BotCount: 2,
		Scenarios: []Scenario{
			{
				Name: "instance_move",
				Parallel: []Scenario{
					{Name: "bot0", Run: func(a *SuiteActor, done func(bool)) { a.InstanceMove(a.bots[0], henesys, done) }},
					{Name: "bot1", Run: func(a *SuiteActor, done func(bool)) { a.InstanceMove(a.bots[1], henesys, done) }},
				},
			},
			{
				Name: "same_map",
				Run: func(a *SuiteActor, done func(bool)) {
					a.Sleep(300*time.Millisecond, func() {
						for _, b := range a.bots {
							if b.Map != henesys {
								done(a.Fail(fmt.Sprintf("%s map %d, want %d", b.Name, b.Map, henesys)))
								return
							}
						}
						done(true)
					})
				},
			},
		},
	},
}
