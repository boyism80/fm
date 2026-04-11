package client

import (
	"net"

	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/types"
)

type Client interface {
	Send(packet types.Packet, policy types.SendPolicy) error
	GetConnection() net.Conn
	GetSendEncryption() *crypt.Encryption
	GetRecvEncryption() *crypt.Encryption
}
