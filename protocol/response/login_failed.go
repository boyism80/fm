package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

const (
	LoginFailedReasonIDDeletedOrBlocked          = 3
	LoginFailedReasonIncorrectPassword           = 4
	LoginFailedReasonNotRegisteredID             = 5
	LoginFailedReasonSystemError6                = 6
	LoginFailedReasonAlreadyLoggedIn             = 7
	LoginFailedReasonSystemError8                = 8
	LoginFailedReasonSystemError9                = 9
	LoginFailedReasonTooManyConnections          = 10
	LoginFailedReasonOnlyAge20AndAbove           = 11
	LoginFailedReasonMasterLoginDenied           = 13
	LoginFailedReasonWrongGatewayOrInfo14        = 14
	LoginFailedReasonKoreanButtonProcessing      = 15
	LoginFailedReasonEmailVerificationRequired16 = 16
	LoginFailedReasonWrongGatewayOrInfo17        = 17
	LoginFailedReasonNoPopup                     = 20
	LoginFailedReasonEmailVerificationRequired21 = 21
	LoginFailedReasonLicenseAgreement            = 23
	LoginFailedReasonMapleEuropeNotice           = 25
	LoginFailedReasonFullClientNotice            = 27
	LoginFailedReasonIPBlocked                   = 32
	LoginFailedReasonPasswordChangeRequired      = 84
)

type LoginFailed struct {
	Reason uint8
}

// Opcode returns the packet opcode for LoginFailed
func (a *LoginFailed) Opcode() uint16 {
	return 0x00
}

func (a *LoginFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(a.Reason)

	switch a.Reason {
	case LoginFailedReasonPasswordChangeRequired:
		writer.WriteDateTime(util.TimeZero)

	case LoginFailedReasonAlreadyLoggedIn:
		writer.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00})
	}

	return nil
}

func (a *LoginFailed) Deserialize(reader *stream.StreamReader) {
}
