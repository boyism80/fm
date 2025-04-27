package resp

import (
	"github.com/boyism80/fm/stream"
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

const (
	TimeFtUtOffset int64 = 116445060000000000 // KST
	TimeMax        int64 = 150842304000000000 // 00 80 05 BB 46 E6 17 02
	TimeZero       int64 = 94354848000000000  // 00 40 E0 FD 3B 37 4F 01
	TimePermanent  int64 = 150841440000000000 // 00 C0 9B 90 7D E5 17 02
)

type LoginFailed struct {
	Reason int
}

func (a *LoginFailed) Serialize(writer *stream.StreamWriter) error {
	err := writer.WriteU16(0x00)
	if err != nil {
		return err
	}

	err = writer.WriteU32(uint32(a.Reason))
	if err != nil {
		return err
	}

	switch a.Reason {
	case LoginFailedReasonPasswordChangeRequired:
		writer.Write64(TimeZero)

	case LoginFailedReasonAlreadyLoggedIn:
		writer.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00})

	}

	return nil
}

func (a *LoginFailed) Deserialize(reader *stream.StreamReader) error {
	return nil
}
