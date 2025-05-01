package msg

import (
	"net"

	"github.com/boyism80/fm/common/context"
)

type StartListening struct {
	Port      int
	ServerCtx *context.ServerContext
}
type ClientConnected struct {
	Conn      net.Conn
	ServerCtx context.IServerContext
}
