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

type CleanupDone struct{}

type SuiteFinished struct {
	Name     string
	Seat     int
	Serial   bool
	Passed   bool
	Skipped  bool
	Infra    bool
	Failures []string
	Elapsed  time.Duration
}
