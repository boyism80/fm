package actor

type Suite struct {
	Name      string
	BotCount  int
	Serial    bool
	Scenarios []Scenario
}

type Scenario struct {
	Name     string
	Run      func(a *SuiteActor, done func(bool))
	Parallel []Scenario
}
