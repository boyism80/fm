package conn

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
)

func TestLoginReceivesAuthenticate(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errc := make(chan error, 1)
	go func() {
		nc, err := ln.Accept()
		if err != nil {
			errc <- err
			return
		}
		defer nc.Close()
		errc <- serveLogin(nc)
	}()

	c, err := Dial(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_ = c.nc.SetDeadline(time.Now().Add(5 * time.Second))

	if err := c.Send(&request.Login{ID: "bot", Pw: "bot-pass", Mac: "AA-BB-CC-DD-EE-FF"}); err != nil {
		t.Fatal(err)
	}
	opcode, body, err := c.Read()
	if err != nil {
		t.Fatal(err)
	}
	pkt, err := Decode(opcode, body)
	if err != nil {
		t.Fatal(err)
	}
	auth, ok := pkt.(*response.Authenticate)
	if !ok {
		t.Fatalf("packet %T", pkt)
	}
	if auth.AccountId != 9 || auth.AccountName != "bot" || auth.Role != 1 {
		t.Fatalf("auth=%+v", auth)
	}
	if err := <-errc; err != nil {
		t.Fatal(err)
	}
}

func serveLogin(nc net.Conn) error {
	_ = nc.SetDeadline(time.Now().Add(5 * time.Second))
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}
	send := crypt.NewEncryption(append([]byte(nil), ivSend...), -5)
	recv := crypt.NewEncryption(append([]byte(nil), ivRecv...), 5)

	welcome := &response.Welcome{SendIv: ivSend, RecvIv: ivRecv}
	w := stream.NewStreamWriter(stream.LittleEndian)
	w.WriteU16(welcome.Opcode())
	if err := welcome.Serialize(w); err != nil {
		return err
	}
	if err := writeFull(nc, w.Bytes()); err != nil {
		return err
	}
	if err := writeEncrypted(nc, &send, &response.LoginFailed{Reason: response.LoginFailedReasonNoPopup}); err != nil {
		return err
	}

	br := bufio.NewReader(nc)
	opcode, body, err := readEncrypted(br, &recv)
	if err != nil {
		return err
	}
	if opcode != uint16((&request.Login{}).Opcode()) {
		return errOpcode(opcode)
	}
	var login request.Login
	login.Deserialize(stream.NewStreamReader(&body, stream.LittleEndian))
	if login.ID != "bot" || login.Pw != "bot-pass" || login.Mac != "AA-BB-CC-DD-EE-FF" {
		return errLogin(login)
	}

	if err := writeEncrypted(nc, &send, &response.Ping{}); err != nil {
		return err
	}
	opcode, _, err = readEncrypted(br, &recv)
	if err != nil {
		return err
	}
	if opcode != uint16((&request.Pong{}).Opcode()) {
		return errOpcode(opcode)
	}
	return writeEncrypted(nc, &send, &response.Authenticate{
		AccountId: 9, Gender: 0, Role: 1, AccountName: "bot",
	})
}

func writeEncrypted(w io.Writer, enc *crypt.Encryption, pkt interface {
	Opcode() uint16
	Serialize(*stream.StreamWriter) error
}) error {
	body := stream.NewStreamWriter(stream.LittleEndian)
	body.WriteU16(pkt.Opcode())
	if err := pkt.Serialize(body); err != nil {
		return err
	}
	plain := body.Bytes()
	frame := append([]byte{}, enc.GetPacketHeader(len(plain))...)
	frame = append(frame, enc.Encrypt(plain)...)
	return writeFull(w, frame)
}

func readEncrypted(r *bufio.Reader, enc *crypt.Encryption) (uint16, []byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}
	if !enc.CheckPacketHeader(header) {
		return 0, nil, errHeader
	}
	n, err := crypt.GetPacketLength(header)
	if err != nil {
		return 0, nil, err
	}
	encBody := make([]byte, n)
	if _, err := io.ReadFull(r, encBody); err != nil {
		return 0, nil, err
	}
	plain := enc.Decrypt(encBody)
	if len(plain) < 2 {
		return 0, nil, errHeader
	}
	opcode := binary.LittleEndian.Uint16(plain[:2])
	body := make([]byte, len(plain)-2)
	copy(body, plain[2:])
	return opcode, body, nil
}

type opcodeError uint16

func (e opcodeError) Error() string { return "opcode" }

func errOpcode(opcode uint16) error { return opcodeError(opcode) }

var errHeader = opcodeError(0)

type loginError struct{ login request.Login }

func (e loginError) Error() string { return "login" }

func errLogin(login request.Login) error { return loginError{login} }
