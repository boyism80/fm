package conn

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
)

type outbound interface {
	Opcode() byte
	Serialize(writer *stream.StreamWriter) error
}

// Conn is one bot TCP session. Send is safe for concurrent callers. Read is not.
type Conn struct {
	nc   net.Conn
	br   *bufio.Reader
	mu   sync.Mutex
	send crypt.Encryption
	recv crypt.Encryption
}

func Dial(addr string) (*Conn, error) {
	nc, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	_ = nc.SetDeadline(time.Now().Add(5 * time.Second))
	c := &Conn{nc: nc, br: bufio.NewReader(nc)}
	if err := c.handshake(); err != nil {
		_ = nc.Close()
		return nil, err
	}
	_ = nc.SetDeadline(time.Time{})
	return c, nil
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.nc.SetDeadline(t)
}

func (c *Conn) Close() error {
	return c.nc.Close()
}

func (c *Conn) handshake() error {
	var hdr [2]byte
	if _, err := io.ReadFull(c.br, hdr[:]); err != nil {
		return err
	}
	n := int(binary.LittleEndian.Uint16(hdr[:]))
	if n <= 0 || n > 256 {
		return fmt.Errorf("welcome length %d", n)
	}
	body := make([]byte, n)
	if _, err := io.ReadFull(c.br, body); err != nil {
		return err
	}
	var welcome response.Welcome
	welcome.Deserialize(stream.NewStreamReader(&body, stream.LittleEndian))
	if len(welcome.RecvIv) != 4 || len(welcome.SendIv) != 4 {
		return fmt.Errorf("welcome iv")
	}
	c.send = crypt.NewEncryption(welcome.RecvIv, 5)
	c.recv = crypt.NewEncryption(welcome.SendIv, -5)
	_, _, err := c.readPacket()
	return err
}

func (c *Conn) Send(pkt outbound) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	w := stream.NewStreamWriter(stream.LittleEndian)
	w.WriteU16(uint16(pkt.Opcode()))
	if err := pkt.Serialize(w); err != nil {
		return err
	}
	plain := w.Bytes()
	frame := append([]byte{}, c.send.GetPacketHeader(len(plain))...)
	frame = append(frame, c.send.Encrypt(plain)...)
	return writeFull(c.nc, frame)
}

// Read returns the next packet. Ping is answered with Pong and not returned.
func (c *Conn) Read() (uint16, []byte, error) {
	for {
		opcode, body, err := c.readPacket()
		if err != nil {
			return 0, nil, err
		}
		if opcode == (&response.Ping{}).Opcode() {
			if err := c.Send(&request.Pong{}); err != nil {
				return 0, nil, err
			}
			continue
		}
		return opcode, body, nil
	}
}

func (c *Conn) readPacket() (uint16, []byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(c.br, header); err != nil {
		return 0, nil, err
	}
	if !c.recv.CheckPacketHeader(header) {
		return 0, nil, fmt.Errorf("invalid packet header")
	}
	n, err := crypt.GetPacketLength(header)
	if err != nil {
		return 0, nil, err
	}
	if n == 0 {
		return 0, nil, fmt.Errorf("empty packet")
	}
	enc := make([]byte, n)
	if _, err := io.ReadFull(c.br, enc); err != nil {
		return 0, nil, err
	}
	plain := c.recv.Decrypt(enc)
	if len(plain) < 2 {
		return 0, nil, fmt.Errorf("packet too short")
	}
	opcode := binary.LittleEndian.Uint16(plain[:2])
	body := make([]byte, len(plain)-2)
	copy(body, plain[2:])
	return opcode, body, nil
}

func writeFull(w io.Writer, b []byte) error {
	for len(b) > 0 {
		n, err := w.Write(b)
		if n > 0 {
			b = b[n:]
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
