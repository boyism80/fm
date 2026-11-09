package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login"
	"github.com/boyism80/fm/types"
)

type DeleteCharacter struct {
	ls *LoginServer
}

func (DeleteCharacter) New(ls *LoginServer) *DeleteCharacter {
	return &DeleteCharacter{
		ls: ls,
	}
}

// deleteRefusals tell the player why internal kept the character.
var deleteRefusals = map[internal.DeleteCharacterResult]string{
	internal.DeleteCharacterResult_DELETE_CHARACTER_RESULT_SECOND_PASSWORD_SET: "2차 비밀번호가 설정된 상태에서는 캐릭터를 삭제할 수 없습니다. 먼저 2차 비밀번호를 해제해주세요.",
	internal.DeleteCharacterResult_DELETE_CHARACTER_RESULT_UNSETTLED:           "처리 중인 거래가 남아 있어 캐릭터를 삭제할 수 없습니다. 게임에 접속해 정리된 뒤 다시 시도해주세요.",
	internal.DeleteCharacterResult_DELETE_CHARACTER_RESULT_SHOP_LEFT:           "프레드릭에게 맡긴 물건이 남아 있어 캐릭터를 삭제할 수 없습니다. 먼저 물건을 찾아가 주세요.",
	internal.DeleteCharacterResult_DELETE_CHARACTER_RESULT_SEALED:              "봉인된 아이템이 있어 캐릭터를 삭제할 수 없습니다.",
}

func (h *DeleteCharacter) Handle(ctx *core.ClientContext, req *request.DeleteCharacter) error {
	log.Printf("Delete character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.ID)

	loginClient := ctx.Client.(*login.LoginClient)
	accountId := loginClient.GetAccountId()
	if accountId == 0 {
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	}

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.DeleteCharacterReply, error) {
			return h.ls.systems.Get[*rpc.InternalConnection]().Client.DeleteCharacter(c, &internal.DeleteCharacterRequest{AccountId: accountId, CharacterId: req.ID, WorldId: loginClient.GetWorldId()})
		}, func(reply *internal.DeleteCharacterReply) error {
			if message, ok := deleteRefusals[reply.GetResult()]; ok {
				_ = ctx.Client.Send(&response.Notice{Type: constant.MsgPopup, Message: message}, types.SEND_POLICY_ENCRYPT)
			}
			deleteResp := &response.DeleteCharacter{
				ID:      req.ID,
				Success: reply.GetResult() == internal.DeleteCharacterResult_DELETE_CHARACTER_RESULT_OK,
			}
			return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
		}).OnError(func(err error) {
		log.Printf("DeleteCharacter (async): %v", err)
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		_ = ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}
