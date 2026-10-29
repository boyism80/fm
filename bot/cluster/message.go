package cluster

const (
	MessageHello         = "hello"
	MessageWelcome       = "welcome"
	MessageRefused       = "refused"
	MessageHeartbeat     = "heartbeat"
	MessageRunSuite      = "run_suite"
	MessageSuiteStarted  = "suite_started"
	MessageSuiteFinished = "suite_finished"
	MessageStop          = "stop"
	MessageQuit          = "quit"
)

// Message is every message between the master and an agent over /ws/agent; Type picks which fields are set.
type Message struct {
	Type string `json:"type"`

	// hello
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty"`
	Host  string `json:"host,omitempty"`
	PID   int    `json:"pid,omitempty"`
	Seats int    `json:"seats,omitempty"`

	// welcome
	AgentID  int    `json:"agent_id,omitempty"`
	SeatBase int    `json:"seat_base,omitempty"`
	Error    string `json:"error,omitempty"`

	// run_suite, suite_started, suite_finished
	Run       string `json:"run,omitempty"`
	Bots      string `json:"bots,omitempty"`
	Suite     string `json:"suite,omitempty"`
	Serial    bool   `json:"serial,omitempty"`
	TimeoutMs int    `json:"timeout_ms,omitempty"`
	Seat      int    `json:"seat,omitempty"`

	// suite_finished
	Title    string   `json:"title,omitempty"`
	State    string   `json:"state,omitempty"`
	Infra    bool     `json:"infra,omitempty"`
	Failures []string `json:"failures,omitempty"`
	Log      []string `json:"log,omitempty"`
	Seconds  float64  `json:"seconds,omitempty"`
}
