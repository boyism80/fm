package msg

import "net"

type StartListening struct {
	Port int
}
type ClientConnected struct {
	Conn net.Conn
}
