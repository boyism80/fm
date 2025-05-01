package resp

import (
	"github.com/boyism80/fm/common/stream"
)

const (
	NoticeTypeDefault                = 0
	NoticeTypePopup                  = 1
	NoticeTypeMegaphone              = 2
	NoticeTypeSuperMegaphone         = 3
	NoticeTypeScrollingMessageTop    = 4
	NoticeTypePinkText               = 5
	NoticeTypeLightblueText          = 6
	NoticeTypeItemMegaphone          = 8
	NoticeTypeHeartMegaphone         = 9
	NoticeTypeSkullSuperMegaphone    = 10
	NoticeTypeGreenMegaphoneMessage  = 11
	NoticeTypeThreeLineMegaphoneText = 12
	NoticeTypeEndOfFile              = 13
	NoticeTypeAniMessage             = 14
	NoticeTypeRedGachaponBox         = 15
	NoticeTypeBlueNoticeAgain        = 18
)

type Notice struct {
	Type    uint8
	Channel uint
	Message string
	MegaEar bool
}

func (a *Notice) Serialize(writer *stream.StreamWriter) error {
	err := writer.WriteU16(0x33)
	if err != nil {
		return err
	}

	err = writer.WriteU8(a.Type)
	if err != nil {
		return err
	}
	if a.Type == 4 {
		writer.WriteU8(1)
	}

	err = writer.WriteStr16(a.Message)
	if err != nil {
		return err
	}

	switch a.Type {
	case NoticeTypeSuperMegaphone:
	case NoticeTypeHeartMegaphone:
	case NoticeTypeSkullSuperMegaphone:
		err = writer.WriteU8(uint8(a.Channel) - 1)
		if err != nil {
			return err
		}
		err = writer.WriteBoolean(a.MegaEar)
		if err != nil {
			return err
		}

	case NoticeTypeLightblueText:
	case NoticeTypeBlueNoticeAgain:
		if a.Channel >= 1000000 && a.Channel < 6000000 {
			err = writer.WriteU32(uint32(a.Channel))
			if err != nil {
				return err
			}
		} else {
			err = writer.WriteU32(0)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Notice) Deserialize(reader *stream.StreamReader) error {
	return nil
}
