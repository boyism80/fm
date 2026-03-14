package server

import (
	"fmt"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
)

type PacketLog struct {
	gs *GameServer
}

func (*PacketLog) New(gs *GameServer) *PacketLog {
	return &PacketLog{
		gs: gs,
	}
}

func (h *PacketLog) GetCommandName() string {
	return "패킷로그"
}

func (h *PacketLog) GetUsage() string {
	return " [0|1|on|off] - 수신 패킷 로그 켜기/끄기 (인자 없으면 현재 상태)"
}

func (h *PacketLog) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) == 0 {
		enabled := core.GetPacketLogEnabled()
		return fmt.Errorf("packet log: %v", enabled)
	}

	arg := strings.ToLower(strings.TrimSpace(args[0]))
	var enable bool
	switch arg {
	case "1", "on", "true":
		enable = true
	case "0", "off", "false":
		enable = false
	default:
		return fmt.Errorf("usage: 0|1|on|off")
	}

	core.SetPacketLogEnabled(enable)
	return fmt.Errorf("packet log set to %v", enable)
}
