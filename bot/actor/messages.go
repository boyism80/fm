package actor

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

type BotEntered struct {
	Bot int
	Err error
}

type ChannelEntered struct {
	Bot int
	Err error
}

type PacketReceived struct {
	Bot    int
	Gen    int
	Opcode uint16
	Body   []byte
}

type Disconnected struct {
	Bot int
	Gen int
	Err error
}

type WaitTimeout struct {
	Wait uint64
}

type SleepDone struct {
	Thread *lua.LState
}

type SuiteDeadline struct{}

type SuiteStop struct{}

type CleanupDone struct{}

type RunStop struct{}

// AssignSuite runs Suite on Seat for a cluster run. Bots replaces the run ID in bot names, so bots of different
// agents in one run do not share a name.
type AssignSuite struct {
	Run   string
	Bots  string
	Suite *Suite
	Seat  int
}

type SuiteStarted struct {
	Suite string
	Seat  int
}

type SuiteFinished struct {
	Suite    string
	Name     string
	Seat     int
	Serial   bool
	Passed   bool
	Skipped  bool
	Stopped  bool
	Infra    bool
	Failures []string
	Log      []string
	Elapsed  time.Duration
}
