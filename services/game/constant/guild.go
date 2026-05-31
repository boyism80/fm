package constant

import "time"

const GuildInviteTargetBusyMessage = "이미 다른 요청을 처리하는 중입니다."

const GuildFullMessage = "가입하려는 길드는 이미 최대 인원으로 가득 찼습니다."

const GuildInviteDuration = 20 * time.Minute

const GuildEmblemMapID uint32 = 200000301

const GuildCreateMesoCost int32 = 1_500_000

const GuildEmblemChangeMesoCost int32 = 5_000_000

const GuildEmblemChangeCashItemID uint32 = 5220001

type GuildCreateResult int

const (
	GuildCreateResultOK               GuildCreateResult = 0
	GuildCreateResultNotAllowed       GuildCreateResult = 1
	GuildCreateResultInsufficientMeso GuildCreateResult = 2
	GuildCreateResultAlreadyInGuild   GuildCreateResult = 3
	GuildCreateResultFailed           GuildCreateResult = 4
	GuildCreateResultSendFailed       GuildCreateResult = -1
)

type GuildDisbandResult int

const (
	GuildDisbandResultOK         GuildDisbandResult = 0
	GuildDisbandResultNotInGuild GuildDisbandResult = 1
	GuildDisbandResultNotMaster  GuildDisbandResult = 2
	GuildDisbandResultFailed     GuildDisbandResult = 3
	GuildDisbandResultSendFailed GuildDisbandResult = -1
)

const GuildEmblemChangeInsufficientCostMessage = "길드마크를 변경할 메소가 부족합니다."
